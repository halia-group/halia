/*
 *
 *  * MIT License
 *  *
 *  * Copyright (c) [2021] [xialeistudio]
 *  *
 *  * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  * of this software and associated documentation files (the "Software"), to deal
 *  * in the Software without restriction, including without limitation the rights
 *  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  * copies of the Software, and to permit persons to whom the Software is
 *  * furnished to do so, subject to the following conditions:
 *  *
 *  * The above copyright notice and this permission notice shall be included in all
 *  * copies or substantial portions of the Software.
 *  *
 *  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  * SOFTWARE.
 *
 */

package common

import (
	"encoding/binary"
	"errors"
	"io"
)

type Packet interface {
	Opcode() uint16
	Write(w io.Writer) error
	Read(r io.Reader) error
}

const (
	_                = iota
	OpPing           = 0x0001
	OpPong           = 0x0002
	OpRegisterReq    = 0x0003
	OpRegisterResult = 0x0004
	OpLoginReq       = 0x0005
	OpLoginResult    = 0x0006
	OpChatReq        = 0x0007
	OpChatResult     = 0x0008
)

var (
	MagicNumber uint16 = 0xcafe
)

var (
	ErrUnknownOpcode = errors.New("unknown opcode")
	ErrInvalidPacket = errors.New("invalid packet")
)

type basePacket struct {
	MagicNumber uint16
	Opcode      uint16
	Length      uint16
	Data        []byte
}

func (p *basePacket) writeCommonField(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, &p.MagicNumber); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, &p.Opcode); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, &p.Length); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, &p.Data); err != nil {
		return err
	}
	return nil
}

func (p *basePacket) readCommonField(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &p.MagicNumber); err != nil {
		return err
	}
	if p.MagicNumber != MagicNumber {
		return ErrInvalidPacket
	}
	if err := binary.Read(r, binary.BigEndian, &p.Opcode); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &p.Length); err != nil {
		return err
	}
	p.Data = make([]byte, p.Length)
	if err := binary.Read(r, binary.BigEndian, &p.Data); err != nil {
		return err
	}
	return nil
}

func (basePacket) writeString(w io.Writer, str string) error {
	if len(str) > 0xff {
		return errors.New("string is too long")
	}
	var (
		buf    = []byte(str)
		length = uint8(len(buf))
	)
	if err := binary.Write(w, binary.BigEndian, &length); err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, &buf)
}

func (basePacket) readString(r io.Reader) (string, error) {
	var (
		length uint8
	)
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return "", err
	}
	buf := make([]byte, length)
	if err := binary.Read(r, binary.BigEndian, &buf); err != nil {
		return "", err
	}
	return string(buf), nil
}
