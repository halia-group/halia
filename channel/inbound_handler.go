package channel

type InboundHandler interface {
	Handler
	ChannelActive(c HandlerContext)
	ChannelInActive(c HandlerContext)
	ChannelRead(c HandlerContext, msg interface{}) (out interface{}, err error)
}
