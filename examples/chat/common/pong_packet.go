package common

import (
	"io"
)

type PongPacket struct {
	basePacket
}

func (p *PongPacket) String() string {
	return "PongPacket{}"
}

func NewPongPacket() *PongPacket {
	return &PongPacket{
		basePacket{
			MagicNumber: MagicNumber,
			Opcode:      OpPong,
			Length:      0,
			Data:        make([]byte, 0),
		},
	}
}

func (p *PongPacket) Opcode() uint16 {
	return OpPong
}

func (p *PongPacket) Write(w io.Writer) error {
	return nil
}

func (p *PongPacket) Read(r io.Reader) error {
	return nil
}
