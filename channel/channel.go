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

package channel

import (
	"bufio"
	"github.com/halia-group/halia/channel/channelid"
	"github.com/halia-group/halia/util"
	"io"
	"net"
)

type Channel interface {
	io.ReadWriteCloser
	util.AttributeMap

	Flush() error
	Id() channelid.ChannelId
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Pipeline() Pipeline
}

type DefaultChannel struct {
	util.DefaultAttributeMap
	conn     net.Conn
	id       channelid.ChannelId
	rw       *bufio.ReadWriter
	pipeline Pipeline
}

func NewDefaultChannel(conn net.Conn) *DefaultChannel {
	c := &DefaultChannel{
		conn: conn,
		id:   channelid.New(),
		rw:   bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
	}
	c.pipeline = NewDefaultPipeline(c)
	return c
}

func (c *DefaultChannel) Read(p []byte) (n int, err error) {
	return c.rw.Read(p)
}

func (c *DefaultChannel) Write(p []byte) (n int, err error) {
	return c.rw.Write(p)
}

func (c *DefaultChannel) Close() error {
	return c.conn.Close()
}

func (c *DefaultChannel) Id() channelid.ChannelId {
	return c.id
}

func (c *DefaultChannel) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *DefaultChannel) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *DefaultChannel) Pipeline() Pipeline {
	return c.pipeline
}

// flush output, only working with buffered writer
func (c *DefaultChannel) Flush() error {
	return c.rw.Flush()
}
