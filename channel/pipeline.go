package channel

type Pipeline interface {
	InboundInvoker
	OutboundInvoker

	AddIn(name string, handler Handler)
	AddOut(name string, handler Handler)
	Channel() Channel
	InNames() []string
	OutNames() []string
}

type DefaultPipeline struct {
	in      *DefaultHandlerContext // 入站链
	out     *DefaultHandlerContext // 出站链
	channel Channel
}

func NewDefaultPipeline(channel Channel) *DefaultPipeline {
	p := &DefaultPipeline{
		channel: channel,
	}
	p.in = &DefaultHandlerContext{
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
	p.out = outHeadCtx
	return p
}

func (p *DefaultPipeline) FireChannelActive() {
	p.in.FireChannelActive()
}

func (p *DefaultPipeline) FireChannelInActive() {
	p.in.FireChannelInActive()
}

func (p *DefaultPipeline) FireChannelRead(msg interface{}) {
	p.in.FireChannelRead(msg)
}

func (p *DefaultPipeline) FireErrorCaught(err error) {
	p.in.FireErrorCaught(err)
}

func (p *DefaultPipeline) Write(msg interface{}) error {
	return p.out.Write(msg)
}

func (p *DefaultPipeline) Flush() error {
	return p.out.Flush()
}

func (p *DefaultPipeline) WriteAndFlush(msg interface{}) error {
	return p.out.WriteAndFlush(msg)
}

// 插入到入站链末尾
func (p *DefaultPipeline) AddIn(name string, handler Handler) {
	newCtx := &DefaultHandlerContext{
		pipeline: p,
		name:     name,
		handler:  handler,
	}
	currentCtx := p.in
	for currentCtx.next != nil {
		currentCtx = currentCtx.next
	}
	currentCtx.next = newCtx
}

// 插入到出站链头部
// 头节点必须是DefaultOutboundHandler, 尾结点必须是TailOutboundHandler
func (p *DefaultPipeline) AddOut(name string, handler Handler) {
	newCtx := &DefaultHandlerContext{
		pipeline: p,
		name:     name,
		handler:  handler,
	}
	newCtx.next = p.out.next
	p.out.next = newCtx
}

func (p *DefaultPipeline) Channel() Channel {
	return p.channel
}

func (p *DefaultPipeline) InNames() []string {
	ptr := p.in
	names := make([]string, 0)
	for ptr != nil {
		names = append(names, ptr.name)
		ptr = ptr.next
	}
	return names
}

func (p *DefaultPipeline) OutNames() []string {
	ptr := p.out
	names := make([]string, 0)
	for ptr != nil {
		names = append(names, ptr.name)
		ptr = ptr.next
	}
	return names
}
