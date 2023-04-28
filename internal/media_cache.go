package internal

import (
	"github.com/bluenviron/gortsplib/v3/pkg/media"
	"github.com/pion/rtp"
)

type mediaCache struct {
	lastMedia  *media.Media
	lastPacket *rtp.Packet
}

func newMediaCache() *mediaCache {
	return &mediaCache{}
}

func (m *mediaCache) cacheReady() bool {
	return m.lastMedia != nil && m.lastPacket != nil
}

func (m *mediaCache) mediaPair() (*media.Media, *rtp.Packet) {
	return m.lastMedia, m.lastPacket
}

func (m *mediaCache) cachePacket(media *media.Media, pkt *rtp.Packet) {
	m.lastMedia = media
	m.lastPacket = pkt
}
