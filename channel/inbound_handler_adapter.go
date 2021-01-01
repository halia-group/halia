package channel

type InboundHandlerAdapter struct{}

func (InboundHandlerAdapter) ErrorCaught(c HandlerContext, err error) {
}

func (InboundHandlerAdapter) ChannelActive(c HandlerContext) {
}

func (InboundHandlerAdapter) ChannelInActive(c HandlerContext) {
}

func (InboundHandlerAdapter) ChannelRead(c HandlerContext, msg interface{}) (out interface{}, err error) {
	return
}
