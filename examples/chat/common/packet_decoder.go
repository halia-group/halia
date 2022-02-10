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
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/handler/codec"
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
