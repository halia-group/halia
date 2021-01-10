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
	client := bootstrap.NewClient(&bootstrap.ClientOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().AddInbound("decoder", &LineDelimiterFrameDecoder{})
			c.Pipeline().AddInbound("handler", NewEchoClientHandler())
			c.Pipeline().AddOutbound("encoder", &StringToByteEncoder{})
			return c
		},
	})

	log.WithField("component", "server").Fatal(client.Dial("tcp", "127.0.0.1:8080"))
}
