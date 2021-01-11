package codec

import "halia/channel"

type Decoder struct{}

func (d Decoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (d Decoder) ChannelActive(c channel.HandlerContext) {
	c.FireChannelActive()
}

func (d Decoder) ChannelInActive(c channel.HandlerContext) {
	c.FireChannelInActive()
}
