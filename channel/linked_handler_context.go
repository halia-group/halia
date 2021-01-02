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

func (c *DefaultHandlerContext) FireErrorCaught(err error) {
	var (
		next    = c.findNextContext()
		handler = next.Handler().(InboundHandler)
	)
	handler.ErrorCaught(next, err)
}

func (c *DefaultHandlerContext) Write(msg interface{}) error {
	var (
		prev    = c.findNextContext()
		handler = prev.Handler().(OutboundHandler)
	)
	return handler.Write(prev, msg)
}

func (c *DefaultHandlerContext) Flush() error {
	var (
		prev    = c.findNextContext()
		handler = prev.Handler().(OutboundHandler)
	)
	return handler.Flush(prev)
}

func (c *DefaultHandlerContext) WriteAndFlush(msg interface{}) error {
	var (
		prev    = c.findNextContext()
		handler = prev.Handler().(OutboundHandler)
	)
	if err := handler.Write(prev, msg); err != nil {
		return err
	}
	return handler.Flush(prev)
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

func (c *DefaultHandlerContext) ErrorCaught(err error) {
	var (
		next    = c.findNextContext()
		handler = next.Handler().(InboundHandler)
	)
	handler.ErrorCaught(next, err)
}

func (c *DefaultHandlerContext) findNextContext() *DefaultHandlerContext {
	return c.next
}
