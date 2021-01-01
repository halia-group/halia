package main

import (
	log "github.com/sirupsen/logrus"
	"halia/bootstrap"
	"halia/channel"
	"net"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	server := bootstrap.NewServer(&bootstrap.ServerOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().Add(NewEchoServerHandler())
			return c
		},
		HandlerContextFactory: func(c channel.Channel) channel.HandlerContext {
			return channel.NewDefaultHandlerContext(c)
		},
	})
	log.Fatal(server.Listen("tcp", "0.0.0.0:8080"))
}
