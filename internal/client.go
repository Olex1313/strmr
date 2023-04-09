package internal

import (
	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/formats"
	"github.com/bluenviron/gortsplib/v3/pkg/media"
	"github.com/bluenviron/gortsplib/v3/pkg/url"
	"github.com/pion/rtp"
	"log"
	"sync"
	"time"
)

const (
	existingStream = "rtsp://192.168.1.123:8554/live"
	reconnectPause = 2 * time.Second
)

type client struct {
	mutex  sync.RWMutex
	stream *gortsplib.ServerStream
}

func NewClient() *client {
	c := &client{}

	// start a separated routine
	go c.run()

	return c
}

func (c *client) run() {
	for {
		err := c.read()
		log.Printf("ERR: %s\n", err)

		time.Sleep(reconnectPause)
	}
}

func (c *client) GetStream() *gortsplib.ServerStream {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stream
}

func (c *client) read() error {
	rc := gortsplib.Client{}

	// parse URL
	u, err := url.Parse(existingStream)
	if err != nil {
		return err
	}

	// connect to the server
	err = rc.Start(u.Scheme, u.Host)
	if err != nil {
		return err // "ECONNREFUSED"
	}
	defer rc.Close()

	// find published medias
	medias, baseURL, _, err := rc.Describe(u)
	if err != nil {
		return err
	}

	// setup all medias
	err = rc.SetupAll(medias, baseURL)
	if err != nil {
		return err // "EOF"
	}

	// create a server stream
	stream := gortsplib.NewServerStream(medias)
	defer stream.Close()

	log.Printf("stream is ready and can be read from the server at rtsp://localhost:8554/stream\n")

	// make stream available by using getStream()
	c.mutex.Lock()
	c.stream = stream
	c.mutex.Unlock()

	defer func() {
		// remove stream from getStream()
		c.mutex.Lock()
		c.stream = nil
		c.mutex.Unlock()
	}()

	// called when RTP packet arrives
	rc.OnPacketRTPAny(func(media *media.Media, format formats.Format, pkt *rtp.Packet) {
		// route incoming packets to the server stream
		stream.WritePacketRTP(media, pkt)
	})

	// start playing
	_, err = rc.Play(nil)
	if err.Error() == "EOF" { // TODO: Write temporary package to stream
		return err
	}
	if err != nil {
		return err // "EOF"
	}

	// wait until a fatal error
	return rc.Wait()
}
