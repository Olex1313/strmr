package internal

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	TcpPort           string `yaml:"tcp-port" envconfig:"SERVER_TCP_PORT" default:":8554"`
	UdpPort           string `yaml:"udp-port" envconfig:"SERVER_UDP_PORT" default:":8000"`
	UdpRtcpPort       string `yaml:"udp-rtcp-port" envconfig:"SERVER_UDP_RTCP_PORT" default:":8001"`
	MulticastIpRange  string `yaml:"multicast-ip-range" envconfig:"SERVER_MULTICAST_IP_RANGE" default:"224.1.0.0/16"`
	MulticastRTPPort  int    `yaml:"multicast-rtp-port" envconfig:"SERVER_MULTICAST_RTP_PORT" default:"8002"`
	MulticastRTCPPort int    `yaml:"multicast-rtcp-port" envconfig:"SERVER_MULTICAST_RTCP_PORT" default:"8003"`
}

type ClientConfig struct {
	ClientAddr     string `yaml:"addr" envconfig:"CLIENT_ADDR"`
	ReconnectPause int    `yaml:"reconnect-pause" envconfig:"CLIENT_RECONNECT_PAUSE"`
}

type ProxyConfig struct {
	ClientConf *ClientConfig `yaml:"client"`
	ServerConf *ServerConfig `yaml:"server"`
}

func NewProxyConfig() *ProxyConfig {
	c := ProxyConfig{}
	readFile("configs/conf.yaml", &c)
	readEnv(&c)
	fmt.Println(c.ClientConf)
	fmt.Println(c.ServerConf)
	return &c
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(path string, cfg interface{}) {
	f, err := os.Open(path)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg interface{}) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
