package common

import (
	"bytes"
	"halia/channel"
	"halia/handler/codec"
)

type PacketDecoder struct {
	codec.Decoder
}

func (p PacketDecoder) ChannelRead(c channel.HandlerContext, msg interface{}) {
	var (
		buf        = msg.([]byte)
		basePacket = basePacket{}
	)
	if err := basePacket.readCommonField(bytes.NewReader(buf)); err != nil {
		c.FireOnError(err)
		return
	}
	packetWrapper, ok := PacketFactory[basePacket.Opcode]
	if !ok {
		c.FireOnError(ErrUnknownOpcode)
		return
	}
	packet := packetWrapper()
	if err := packet.Read(bytes.NewReader(basePacket.Data)); err != nil {
		c.FireOnError(err)
		return
	}
	c.FireChannelRead(packet)
}
