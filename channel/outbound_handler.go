package channel

type OutboundHandler interface {
	Handler
	Write(c HandlerContext, msg interface{}) (out interface{}, err error)
	Flush(c HandlerContext) error
}
