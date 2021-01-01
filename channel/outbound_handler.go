package channel

type OutboundHandler interface {
	Handler
	Write(c HandlerContext, msg interface{}) error
	Flush(c HandlerContext) error
}
