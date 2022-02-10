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
	"fmt"
	"io"
)

type RegisterReq struct {
	basePacket
	Username string
	Password string
}

func (p *RegisterReq) String() string {
	return fmt.Sprintf("RegisterReq{Username=%s,Password=%s}", p.Username, p.Password)
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

func (p *RegisterResult) String() string {
	return fmt.Sprintf("RegisterResult{Code=%d,Message=%s}", p.Code, p.Message)
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
