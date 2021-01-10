package channel

type Pipeline interface {
	InboundInvoker
	OutboundInvoker

	AddInbound(name string, handler Handler)
	AddOutbound(name string, handler Handler)
	Channel() Channel
	InboundNames() []string
	OutboundNames() []string
}

type DefaultPipeline struct {
	inbound  *DefaultHandlerContext // 入站链
	outbound *DefaultHandlerContext // 出站链
	channel  Channel
}

func NewDefaultPipeline(channel Channel) *DefaultPipeline {
	p := &DefaultPipeline{
		channel: channel,
	}
	p.inbound = &DefaultHandlerContext{
		pipeline: p,
		name:     "InHeadContext",
		handler:  NewHeadInboundHandler(p),
	}
	outTailCtx := &DefaultHandlerContext{
		pipeline: p,
		name:     "OutTailContext",
		handler:  NewTailOutboundHandler(p),
	}
	outHeadCtx := &DefaultHandlerContext{
		next:     outTailCtx,
		pipeline: p,
		name:     "OutHeadContext",
		handler:  NewHeadOutboundHandler(p),
	}
	p.outbound = outHeadCtx
	return p
}

func (p *DefaultPipeline) FireChannelActive() {
	p.inbound.FireChannelActive()
}

func (p *DefaultPipeline) FireChannelInActive() {
	p.inbound.FireChannelInActive()
}

func (p *DefaultPipeline) FireChannelRead(msg interface{}) {
	p.inbound.FireChannelRead(msg)
}

func (p *DefaultPipeline) FireOnError(err error) {
	p.inbound.FireOnError(err)
}

func (p *DefaultPipeline) Write(msg interface{}) error {
	return p.outbound.Write(msg)
}

func (p *DefaultPipeline) Flush() error {
	return p.outbound.Flush()
}

func (p *DefaultPipeline) WriteAndFlush(msg interface{}) error {
	return p.outbound.WriteAndFlush(msg)
}

// 插入到入站链末尾
func (p *DefaultPipeline) AddInbound(name string, handler Handler) {
	newCtx := &DefaultHandlerContext{
		pipeline: p,
		name:     name,
		handler:  handler,
	}
	currentCtx := p.inbound
	for currentCtx.next != nil {
		currentCtx = currentCtx.next
	}
	currentCtx.next = newCtx
}

// 插入到出站链头部
// 头节点必须是DefaultOutboundHandler, 尾结点必须是TailOutboundHandler
func (p *DefaultPipeline) AddOutbound(name string, handler Handler) {
	newCtx := &DefaultHandlerContext{
		pipeline: p,
		name:     name,
		handler:  handler,
	}
	newCtx.next = p.outbound.next
	p.outbound.next = newCtx
}

func (p *DefaultPipeline) Channel() Channel {
	return p.channel
}

func (p *DefaultPipeline) InboundNames() []string {
	ptr := p.inbound
	names := make([]string, 0)
	for ptr != nil {
		names = append(names, ptr.name)
		ptr = ptr.next
	}
	return names
}

func (p *DefaultPipeline) OutboundNames() []string {
	ptr := p.outbound
	names := make([]string, 0)
	for ptr != nil {
		names = append(names, ptr.name)
		ptr = ptr.next
	}
	return names
}
