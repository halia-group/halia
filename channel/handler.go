package channel

type Handler interface {
	ErrorCaught(c HandlerContext, err error)
}
