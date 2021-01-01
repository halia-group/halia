package bootstrap

import (
	"halia/channel"
	"net"
)

type ServerOptions struct {
	ChannelFactory        func(net.Conn) channel.Channel
	HandlerContextFactory func(channel2 channel.Channel) channel.HandlerContext
}
