package channel

type OutboundHandlerAdapter struct{}

func (OutboundHandlerAdapter) ErrorCaught(c HandlerContext, err error) {
}

func (OutboundHandlerAdapter) Write(c HandlerContext, msg interface{}) (out interface{}, err error) {
	return
}

func (OutboundHandlerAdapter) Flush(c HandlerContext) error {
	return nil
}
