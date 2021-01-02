package bootstrap

import (
	"halia/channel"
	"net"
)

type ServerOptions struct {
	ChannelFactory func(conn net.Conn) channel.Channel
}
