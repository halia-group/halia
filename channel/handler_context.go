package channel

type HandlerContext interface {
	InboundInvoker
	OutboundInvoker

	Channel() Channel
	Name() string
	Handler() Handler
	Pipeline() Pipeline
	OnError(err error)
}
