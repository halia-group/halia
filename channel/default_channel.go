package channel

import (
	"halia/channel/channelid"
	"net"
)

type DefaultChannel struct {
	id       channelid.ChannelId
	pipeline Pipeline
	conn     net.Conn
}

func NewDefaultChannel(conn net.Conn) *DefaultChannel {
	return &DefaultChannel{
		id:       channelid.New(),
		conn:     conn,
		pipeline: NewDefaultPipeline(),
	}
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

func (c *DefaultChannel) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *DefaultChannel) WriteAndFlush(data []byte) error {
	if _, err := c.Write(data); err != nil {
		return err
	}
	if err := c.Flush(); err != nil {
		return err
	}
	return nil
}

func (c *DefaultChannel) Flush() error {
	return nil
}

func (c *DefaultChannel) Close() error {
	return c.conn.Close()
}

func (c *DefaultChannel) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}
