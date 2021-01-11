package common

import (
	"bytes"
	"encoding/binary"
	"halia/channel"
)

type PacketEncoder struct{}

func (p PacketEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (p PacketEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	packet := msg.(Packet)
	buf := bytes.Buffer{}
	if err := packet.Write(&buf); err != nil {
		return err
	}
	var (
		body   = buf.Bytes()
		opcode = packet.Opcode()
		length = uint16(len(body))
	)
	buf = bytes.Buffer{}
	if err := binary.Write(&buf, binary.BigEndian, &MagicNumber); err != nil {
		return err
	}
	if err := binary.Write(&buf, binary.BigEndian, &opcode); err != nil {
		return err
	}
	if err := binary.Write(&buf, binary.BigEndian, &length); err != nil {
		return err
	}
	if err := binary.Write(&buf, binary.BigEndian, &body); err != nil {
		return err
	}
	return c.Write(buf.Bytes())
}

func (p PacketEncoder) Flush(c channel.HandlerContext) error {
	return c.Channel().Flush()
}
