package main

import (
	"flag"
	"gobs/internal"
	"time"
)

func main() {
	configPath := flag.String("config", "configs/conf.yaml", "Path to the config file")
	flag.Parse()
	conf := internal.NewProxyConfig(configPath)
	c := internal.NewClient(
		conf.ClientConf.ClientAddr,
		time.Millisecond*time.Duration(conf.ClientConf.ReconnectPause),
		time.Millisecond*time.Duration(conf.ClientConf.ReconnectInterval),
	)
	internal.NewServer(conf.ServerConf, c.GetStream)
}
