package channel

type Pipeline interface {
	InboundInvoker
	OutboundInvoker

	AddFirst(name string, handler Handler)
	AddLast(name string, handler Handler)
	AddBefore(baseName, name string, handler Handler)
	AddAfter(baseName, name string, handler Handler)
	Remove(name string)
	First() Handler
	FirstContext() HandlerContext
	Last() Handler
	LastContext() HandlerContext
	Get(name string) Handler
	Context(name string) HandlerContext
	Channel() Channel
	Names() []string
}
