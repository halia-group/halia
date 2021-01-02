package channel

import (
	"halia/channel/channelid"
	"halia/util"
	"io"
	"net"
)

type Channel interface {
	io.ReadWriteCloser
	util.AttributeMap

	Id() channelid.ChannelId
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Pipeline() Pipeline
}
