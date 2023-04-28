package internal

import (
	"context"
	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/media"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
	"log"
	"sync"
	"time"
)

type client struct {
	mutex    sync.RWMutex
	stream   *gortsplib.ServerStream
	mCache   *mediaCache
	addr     string
	recPause time.Duration
}

func NewClient(addr string, recp time.Duration) *client {
	c := &client{
		addr:     addr,
		recPause: recp,
	}
	c.mCache = newMediaCache()

	// start a separated routine
	go c.run()

	return c
}

func (c *client) run() {
	for {
		err := c.read()
		if err != nil {
			break
		}
		log.Printf("ERR: %s\n", err)
	}
}

func (c *client) close() {
	func() {
		// remove stream from getStream()
		c.mutex.Lock()
		c.stream = nil
		c.mutex.Unlock()
	}()
}

func (c *client) GetStream() *gortsplib.ServerStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stream
}

func (c *client) connect() (*gortsplib.Client, error) {
	rc := gortsplib.Client{}

	// parse URL
	u, err := url.Parse(c.addr)
	if err != nil {
		return nil, err
	}

	// connect to the server
	err = rc.Start(u.Scheme, u.Host)
	if err != nil {
		return nil, err // "ECONNREFUSED"
	}

	// find published medias
	medias, baseURL, _, err := rc.Describe(u)
	if err != nil {
		return nil, err
	}

	// setup all medias
	err = rc.SetupAll(medias, baseURL)
	if err != nil {
		return nil, err // "EOF"
	}

	// create a server stream
	stream := gortsplib.NewServerStream(medias)

	log.Printf("stream is ready and can be read from the server at rtsp://localhost:8554/stream\n")

	// make stream available by using getStream()
	c.mutex.Lock()
	c.stream = stream
	c.mutex.Unlock()

	// called when RTP packet arrives
	rc.OnPacketRTPAny(func(media *media.Media, format formats.Format, pkt *rtp.Packet) {
		c.mCache.cachePacket(media, pkt)
		stream.WritePacketRTP(media, pkt)
	})

	// start playing
	_, err = rc.Play(nil)
	return &rc, nil
}

func (c *client) awaitForReconnect(ctx context.Context, ch chan bool) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Failed to reconnect")
			if ctx.Err() == context.DeadlineExceeded {
				panic("TimedOut")
			}
			return
		case <-ch:
		}
	}
}

func (c *client) publishCache(ctx context.Context, ch chan bool) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Resuming stream")
			return
		case <-ch:
			c.stream.WritePacketRTP(c.mCache.mediaPair())
		}
	}
}

func (c *client) read() error {
	rc, err := c.connect()
	// wait until a fatal error or client fails
	err = rc.Wait()
	if err.Error() == "EOF" && c.mCache.cacheReady() {
		ctx, cancel := context.WithTimeout(context.Background(), c.recPause)
		ch := make(chan bool)
		go c.awaitForReconnect(ctx, ch)
		go c.publishCache(ctx, ch)
		rc = c.tryReconnect()
		cancel()
	} else if err != nil {
		return err
	}
	return nil
}

func (c *client) tryReconnect() *gortsplib.Client {
	for {
		rc, err := c.connect()
		if err != nil {
			log.Printf("Reconnect try failed... %s", err.Error())
			continue
		}
		return rc
	}
}
