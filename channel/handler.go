package channel

type Handler interface {
	ErrorCaught(c HandlerContext, err error)
}

type InboundHandler interface {
	Handler
	ChannelActive(c HandlerContext)
	ChannelInActive(c HandlerContext)
	ChannelRead(c HandlerContext, msg interface{})
}

type OutboundHandler interface {
	Handler
	Write(c HandlerContext, msg interface{}) error
	Flush(c HandlerContext) error
}
