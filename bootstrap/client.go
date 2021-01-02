package bootstrap

import (
	"halia/channel"
	"net"
)

type Client struct {
	options *ClientOptions
}

func NewClient(options *ClientOptions) *Client {
	return &Client{
		options: options,
	}
}

func (client *Client) Dial(network, addr string) error {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return err
	}
	c := client.options.ChannelFactory(conn)

	defer c.Pipeline().FireChannelInActive()

	c.Pipeline().FireChannelActive()

	c.Pipeline().FireChannelRead(nil)

	return nil
}

type ClientOptions struct {
	ChannelFactory func(net.Conn) channel.Channel
}
