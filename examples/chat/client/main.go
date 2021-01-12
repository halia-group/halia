package main

import (
	"encoding/binary"
	"github.com/halia-group/halia/bootstrap"
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/examples/chat/common"
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
			c.Pipeline().AddInbound("frameDecoder", codec.NewLengthFieldBasedFrameDecoder(2, 4, binary.BigEndian))
			c.Pipeline().AddInbound("packetDecoder", &common.PacketDecoder{})
			c.Pipeline().AddInbound("handler", newChatClientHandler())

			c.Pipeline().AddOutbound("packetEncoder", &common.PacketEncoder{})
			return c
		},
	})

	log.WithField("component", "client").Fatal(client.Dial("tcp", "127.0.0.1:8080"))
}
