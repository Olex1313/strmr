package main

import (
	"gobs/internal"
	"time"
)

func main() {
	conf := internal.NewProxyConfig()
	c := internal.NewClient(conf.ClientConf.ClientAddr, time.Millisecond*time.Duration(conf.ClientConf.ReconnectPause))
	internal.NewServer(conf.ServerConf, c.GetStream)
}
