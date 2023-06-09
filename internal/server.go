package internal

import (
	"github.com/bluenviron/gortsplib/v3"
	"github.com/bluenviron/gortsplib/v3/pkg/base"
	"log"
)

type server struct {
	getStream func() *gortsplib.ServerStream
}

func NewServer(
	serverConf *ServerConfig,
	getStream func() *gortsplib.ServerStream,
) *server {
	s := &server{
		getStream: getStream,
	}

	// configure the server
	rs := &gortsplib.Server{
		Handler:           s,
		RTSPAddress:       serverConf.TcpPort,
		UDPRTPAddress:     serverConf.UdpPort,
		UDPRTCPAddress:    serverConf.UdpRtcpPort,
		MulticastIPRange:  serverConf.MulticastIpRange,
		MulticastRTPPort:  serverConf.MulticastRTPPort,
		MulticastRTCPPort: serverConf.MulticastRTCPPort,
	}

	// start server and wait until a fatal error
	log.Printf("server is ready")
	panic(rs.StartAndWait())
}

// OnConnOpen called when a connection is opened.
func (s *server) OnConnOpen(ctx *gortsplib.ServerHandlerOnConnOpenCtx) {
	log.Printf("conn opened")
}

// OnConnClose called when a connection is closed.
func (s *server) OnConnClose(ctx *gortsplib.ServerHandlerOnConnCloseCtx) {
	log.Printf("conn closed (%v)", ctx.Error)
}

// OnSessionOpen called when a session is opened.
func (s *server) OnSessionOpen(ctx *gortsplib.ServerHandlerOnSessionOpenCtx) {
	log.Printf("session opened")
}

// OnSessionClose called when a session is closed.
func (s *server) OnSessionClose(ctx *gortsplib.ServerHandlerOnSessionCloseCtx) {
	log.Printf("session closed")
}

// OnDescribe called when receiving a DESCRIBE request.
func (s *server) OnDescribe(ctx *gortsplib.ServerHandlerOnDescribeCtx) (*base.Response, *gortsplib.ServerStream, error) {
	log.Printf("describe request")

	stream := s.getStream()

	// stream is not available yet
	if stream == nil {
		return &base.Response{
			StatusCode: base.StatusNotFound,
		}, nil, nil
	}

	return &base.Response{
		StatusCode: base.StatusOK,
	}, stream, nil
}

// OnSetup called when receiving a SETUP request.
func (s *server) OnSetup(ctx *gortsplib.ServerHandlerOnSetupCtx) (*base.Response, *gortsplib.ServerStream, error) {
	log.Printf("setup request")

	stream := s.getStream()

	// stream is not available yet
	if stream == nil {
		return &base.Response{
			StatusCode: base.StatusNotFound,
		}, nil, nil
	}

	return &base.Response{
		StatusCode: base.StatusOK,
	}, stream, nil
}

// OnPlay called when receiving a PLAY request.
func (s *server) OnPlay(ctx *gortsplib.ServerHandlerOnPlayCtx) (*base.Response, error) {
	log.Printf("play request")

	return &base.Response{
		StatusCode: base.StatusOK,
	}, nil
}
