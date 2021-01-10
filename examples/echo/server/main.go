package main

import (
	log "github.com/sirupsen/logrus"
	"halia/bootstrap"
	"halia/channel"
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
			c.Pipeline().AddInbound("decoder", &LineDelimiterFrameDecoder{})
			c.Pipeline().AddInbound("handler", NewEchoServerHandler())
			c.Pipeline().AddOutbound("encoder", &StringToByteEncoder{})
			return c
		},
	})

	log.WithField("component", "server").Fatal(s.Listen("tcp", "0.0.0.0:8080"))
}
