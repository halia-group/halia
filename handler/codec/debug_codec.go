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

package codec

import (
	"encoding/hex"
	"fmt"
	"github.com/halia-group/halia/channel"
	"os"
)

type DebugDecoder struct {
	Decoder
}

func NewDebugDecoder() *DebugDecoder {
	return &DebugDecoder{}
}

func (d DebugDecoder) ChannelRead(c channel.HandlerContext, msg interface{}) {
	buf := msg.([]byte)
	if _, err := hex.Dumper(os.Stderr).Write(buf); err != nil {
		d.OnError(c, err)
	}
	fmt.Fprintln(os.Stderr)
	c.FireChannelRead(msg)
}

type DebugEncoder struct{}

func NewDebugEncoder() *DebugEncoder {
	return &DebugEncoder{}
}

func (d DebugEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (d DebugEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	buf := msg.([]byte)
	if _, err := hex.Dumper(os.Stderr).Write(buf); err != nil {
		d.OnError(c, err)
	}
	fmt.Fprintln(os.Stderr)
	return c.Write(msg)
}

func (d DebugEncoder) Flush(c channel.HandlerContext) error {
	return c.Flush()
}
