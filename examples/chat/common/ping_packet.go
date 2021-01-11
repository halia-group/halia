package common

import (
	"io"
)

type PingPacket struct {
	basePacket
}

func NewPingPacket() *PingPacket {
	return &PingPacket{
		basePacket{
			MagicNumber: MagicNumber,
			Opcode:      OpPing,
			Length:      0,
			Data:        make([]byte, 0),
		},
	}
}

func (p *PingPacket) Opcode() uint16 {
	return OpPing
}

func (p *PingPacket) Write(w io.Writer) error {
	return nil
}

func (p *PingPacket) Read(r io.Reader) error {
	return nil
}
