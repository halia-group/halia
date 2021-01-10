package codec

import "halia/channel"

type FrameDecoder struct{}

func (d FrameDecoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (d FrameDecoder) ChannelActive(c channel.HandlerContext) {
	c.FireChannelActive()
}

func (d FrameDecoder) ChannelInActive(c channel.HandlerContext) {
	c.FireChannelInActive()
}
