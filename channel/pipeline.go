package channel

type Pipeline interface {
	Add(handler Handler)
	IterateInbound(func(handler InboundHandler) error) error
	IterateOutbound(func(handler OutboundHandler) error) error
}
