package channel

type Pipeline interface {
	InboundInvoker
	OutboundInvoker

	AddIn(name string, handler Handler)
	AddOut(name string, handler Handler)
	First() Handler
	FirstContext() HandlerContext
	Last() Handler
	LastContext() HandlerContext
	Channel() Channel
	Names() []string
}
