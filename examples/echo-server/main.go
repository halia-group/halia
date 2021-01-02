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
	server := bootstrap.NewServer(&bootstrap.ServerOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().AddIn("decoder", &LineDelimiterFrameDecoder{})
			c.Pipeline().AddIn("handler", NewEchoServerHandler())
			c.Pipeline().AddOut("encoder", &StringToByteEncoder{})
			return c
		},
	})

	log.WithField("component", "server").Fatal(server.Listen("tcp", "0.0.0.0:8080"))
}
