package channel

type DefaultHandlerContext struct {
	next     *DefaultHandlerContext
	pipeline *DefaultPipeline
	name     string
	handler  Handler
}

func (c *DefaultHandlerContext) FireChannelActive() {
	var (
		next    = c.findNextContext()
		handler = next.Handler().(InboundHandler)
	)
	handler.ChannelActive(next)
}

func (c *DefaultHandlerContext) FireChannelInActive() {
	var (
		next    = c.findNextContext()
		handler = next.Handler().(InboundHandler)
	)
	handler.ChannelInActive(next)
}

func (c *DefaultHandlerContext) FireChannelRead(msg interface{}) {
	var (
		next    = c.findNextContext()
		handler = next.Handler().(InboundHandler)
	)
	handler.ChannelRead(next, msg)
}

func (c *DefaultHandlerContext) FireOnError(err error) {
	var (
		next    = c.findNextContext()
		handler = next.Handler().(InboundHandler)
	)
	handler.OnError(next, err)
}

// 如果当前是InboundHandler，则调用pipeline的write开始出站流程
func (c *DefaultHandlerContext) Write(msg interface{}) error {
	if _, ok := c.handler.(InboundHandler); ok {
		return c.pipeline.Write(msg)
	}

	var (
		next    = c.findNextContext()
		handler = next.Handler().(OutboundHandler)
	)
	return handler.Write(next, msg)
}

func (c *DefaultHandlerContext) Flush() error {
	if _, ok := c.handler.(InboundHandler); ok {
		return c.pipeline.Flush()
	}

	var (
		next    = c.findNextContext()
		handler = next.Handler().(OutboundHandler)
	)
	return handler.Flush(next)
}

func (c *DefaultHandlerContext) WriteAndFlush(msg interface{}) error {
	if _, ok := c.handler.(InboundHandler); ok {
		return c.pipeline.WriteAndFlush(msg)
	}

	var (
		next    = c.findNextContext()
		handler = next.Handler().(OutboundHandler)
	)
	if err := handler.Write(next, msg); err != nil {
		return err
	}
	return handler.Flush(next)
}

func (c *DefaultHandlerContext) Channel() Channel {
	return c.pipeline.channel
}

func (c *DefaultHandlerContext) Name() string {
	return c.name
}

func (c *DefaultHandlerContext) Handler() Handler {
	return c.handler
}

func (c *DefaultHandlerContext) Pipeline() Pipeline {
	return c.pipeline
}

func (c *DefaultHandlerContext) OnError(err error) {
	var (
		next    = c.findNextContext()
		handler = next.Handler().(InboundHandler)
	)
	handler.OnError(next, err)
}

func (c *DefaultHandlerContext) findNextContext() *DefaultHandlerContext {
	return c.next
}
