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
