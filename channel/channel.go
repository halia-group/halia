package channel

import (
	"halia/channel/channelid"
	"net"
)

type Channel interface {
	OutboundInvoker
	Id() channelid.ChannelId
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Pipeline() Pipeline
}
