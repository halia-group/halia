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
	client := bootstrap.NewClient(&bootstrap.ClientOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().AddLast("decoder", codec.NewLineBasedFrameDecoder())
			c.Pipeline().AddLast("encoder", &StringToByteEncoder{})
			c.Pipeline().AddLast("handler", NewEchoClientHandler())
			return c
		},
	})

	log.WithField("component", "server").Fatal(client.Dial("tcp", "127.0.0.1:8080"))
}
