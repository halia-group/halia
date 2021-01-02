package channel

import log "github.com/sirupsen/logrus"

type Pipeline interface {
	InboundInvoker
	OutboundInvoker

	AddIn(name string, handler Handler)
	AddOut(name string, handler Handler)
	Channel() Channel
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
		name:     "HeadContext",
		handler:  NewHeadInboundHandler(p),
	}
	p.out = &DefaultHandlerContext{
		pipeline: p,
		name:     "TailContext",
		handler:  NewTailOutboundHandler(p),
	}
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
func (p *DefaultPipeline) AddOut(name string, handler Handler) {
	newCtx := &DefaultHandlerContext{
		pipeline: p,
		name:     name,
		handler:  handler,
	}
	newCtx.next = p.out
	p.out = newCtx
}

func (p *DefaultPipeline) Channel() Channel {
	return p.channel
}

// 第一个入站处理器
type HeadInboundHandler struct {
	pipeline *DefaultPipeline
}

func NewHeadInboundHandler(pipeline *DefaultPipeline) *HeadInboundHandler {
	return &HeadInboundHandler{pipeline: pipeline}
}

func (p *HeadInboundHandler) ErrorCaught(c HandlerContext, err error) {
	c.FireErrorCaught(err)
}

func (p *HeadInboundHandler) ChannelActive(c HandlerContext) {
	c.FireChannelActive()
}

func (p *HeadInboundHandler) ChannelInActive(c HandlerContext) {
	c.FireChannelInActive()
}

func (p *HeadInboundHandler) ChannelRead(c HandlerContext, msg interface{}) {
	c.FireChannelRead(msg)
}

// 最后一个出站处理器
type TailOutboundHandler struct {
	pipeline *DefaultPipeline
	log      *log.Entry
}

func NewTailOutboundHandler(pipeline *DefaultPipeline) *TailOutboundHandler {
	return &TailOutboundHandler{pipeline: pipeline, log: log.WithField("component", "TailOutboundHandler")}
}

func (p *TailOutboundHandler) ErrorCaught(c HandlerContext, err error) {
	p.log.Warnln("unhandled error.", err)
}

func (p *TailOutboundHandler) Write(c HandlerContext, msg interface{}) error {
	buf, ok := msg.([]byte)
	if !ok {
		p.log.Warnln("TailOutboundHandler.Write called with wrong msg type")
		return nil
	}
	if _, err := c.Channel().Write(buf); err != nil {
		return err
	}
	return nil
}

func (p *TailOutboundHandler) Flush(c HandlerContext) error {
	return c.Channel().Flush()
}
