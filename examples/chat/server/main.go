package main

import (
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	"halia/bootstrap"
	"halia/channel"
	"halia/examples/chat/common"
	"halia/handler/codec"
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
			c.Pipeline().AddInbound("frameDecoder", codec.NewLengthFieldBasedFrameDecoder(2, 4, binary.BigEndian))
			c.Pipeline().AddInbound("packetDecoder", &common.PacketDecoder{})
			c.Pipeline().AddInbound("handler", newChatServerHandler())

			c.Pipeline().AddOutbound("packetEncoder", &common.PacketEncoder{})
			return c
		},
	})

	log.WithField("component", "server").Fatal(s.Listen("tcp", "0.0.0.0:8080"))
}
