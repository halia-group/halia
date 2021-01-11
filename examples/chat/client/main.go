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
	client := bootstrap.NewClient(&bootstrap.ClientOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().AddInbound("frameDecoder", codec.NewLengthFieldBasedFrameDecoder(2, 4, binary.BigEndian))
			c.Pipeline().AddInbound("debugDecoder", codec.NewDebugDecoder())
			c.Pipeline().AddInbound("packetDecoder", &common.PacketDecoder{})
			c.Pipeline().AddInbound("handler", newChatClientHandler())

			c.Pipeline().AddOutbound("debugEncoder", codec.NewDebugEncoder())
			c.Pipeline().AddOutbound("packetEncoder", &common.PacketEncoder{})
			return c
		},
	})

	log.WithField("component", "client").Fatal(client.Dial("tcp", "127.0.0.1:8080"))
}
