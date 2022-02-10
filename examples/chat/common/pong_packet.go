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
