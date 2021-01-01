package channel

type DefaultPipeline struct {
	inbounds  []InboundHandler
	outbounds []OutboundHandler
}

func NewDefaultPipeline() *DefaultPipeline {
	return &DefaultPipeline{
		inbounds:  make([]InboundHandler, 0),
		outbounds: make([]OutboundHandler, 0),
	}
}

func (p *DefaultPipeline) Add(handler Handler) {
	if h, ok := handler.(InboundHandler); ok {
		p.inbounds = append(p.inbounds, h)
		return
	}
	if h, ok := handler.(OutboundHandler); ok {
		p.outbounds = append(p.outbounds, h)
	}
}

func (p *DefaultPipeline) IterateInbound(f func(handler InboundHandler) error) error {
	for i := range p.inbounds {
		if err := f(p.inbounds[i]); err != nil {
			return err
		}
	}
	return nil
}

func (p *DefaultPipeline) IterateOutbound(f func(handler OutboundHandler) error) error {
	var length = len(p.outbounds)

	for i := length - 1; i > 0; i-- {
		if err := f(p.outbounds[i]); err != nil {
			return err
		}
	}
	return nil
}
