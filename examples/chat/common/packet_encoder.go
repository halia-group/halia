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
	"bytes"
	"encoding/binary"
	"github.com/halia-group/halia/channel"
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
	if err := binary.Write(c.Channel(), binary.BigEndian, &MagicNumber); err != nil {
		return err
	}
	if err := binary.Write(c.Channel(), binary.BigEndian, &opcode); err != nil {
		return err
	}
	if err := binary.Write(c.Channel(), binary.BigEndian, &length); err != nil {
		return err
	}
	if err := binary.Write(c.Channel(), binary.BigEndian, &body); err != nil {
		return err
	}
	return c.Flush()
}

func (p PacketEncoder) Flush(c channel.HandlerContext) error {
	return c.Channel().Flush()
}
