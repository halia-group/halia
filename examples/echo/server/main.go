package main

import (
	"github.com/halia-group/halia/bootstrap"
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/handler/codec"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	s := bootstrap.NewServer(&bootstrap.ServerOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().AddLast("decoder", codec.NewLineBasedFrameDecoder())
			c.Pipeline().AddLast("encoder", &StringToByteEncoder{})
			c.Pipeline().AddLast("handler", NewEchoServerHandler())
			return c
		},
	})

	log.WithField("component", "server").Fatal(s.Listen("tcp", "0.0.0.0:8080"))
}
