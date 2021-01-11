package common

import (
	"encoding/binary"
	"io"
)

type RegisterReq struct {
	basePacket
	Username string
	Password string
}

func (p *RegisterReq) Opcode() uint16 {
	return OpRegisterReq
}

func (p *RegisterReq) Write(w io.Writer) error {
	if err := p.writeString(w, p.Username); err != nil {
		return err
	}
	if err := p.writeString(w, p.Password); err != nil {
		return err
	}
	return nil
}

func (p *RegisterReq) Read(r io.Reader) (err error) {
	if p.Username, err = p.readString(r); err != nil {
		return
	}
	if p.Password, err = p.readString(r); err != nil {
		return
	}
	return
}

func NewRegisterReq(username string, password string) *RegisterReq {
	return &RegisterReq{
		basePacket: basePacket{
			MagicNumber: MagicNumber,
			Opcode:      OpRegisterReq,
		},
		Username: username,
		Password: password,
	}
}

type RegisterResult struct {
	basePacket
	Code    uint8
	Message string
}

func NewRegisterResult(code uint8, message string) *RegisterResult {
	return &RegisterResult{
		basePacket: basePacket{
			MagicNumber: MagicNumber,
			Opcode:      OpRegisterResult,
		},
		Code:    code,
		Message: message,
	}
}

func (p *RegisterResult) Opcode() uint16 {
	return OpRegisterResult
}

func (p *RegisterResult) Write(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, &p.Code); err != nil {
		return err
	}
	if err := p.writeString(w, p.Message); err != nil {
		return err
	}
	return nil
}

func (p *RegisterResult) Read(r io.Reader) (err error) {
	if err = binary.Read(r, binary.BigEndian, &p.Code); err != nil {
		return
	}
	if p.Message, err = p.readString(r); err != nil {
		return
	}
	return
}
