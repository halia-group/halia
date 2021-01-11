package common

import (
	"io"
)

type PingPacket struct {
	basePacket
}

func (p *PingPacket) String() string {
	return "PingPacket{}"
}

func NewPingPacket() *PingPacket {
	return &PingPacket{
		basePacket{
			MagicNumber: MagicNumber,
			Opcode:      OpPing,
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
