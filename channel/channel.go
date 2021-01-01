package channel

import (
	"halia/channel/channelid"
	"io"
	"net"
)

type Channel interface {
	io.ReadWriteCloser
	Id() channelid.ChannelId
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Pipeline() Pipeline
	WriteAndFlush(data []byte) error
	Flush() error
}
