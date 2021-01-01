package channel

import "halia/util"

type DefaultHandlerContext struct {
	channel Channel
	attrMap util.AttributeMap
}

func NewDefaultHandlerContext(channel Channel) *DefaultHandlerContext {
	return &DefaultHandlerContext{channel: channel, attrMap: util.NewDefaultAttributeMap()}
}

func (c *DefaultHandlerContext) Channel() Channel {
	return c.channel
}

func (c *DefaultHandlerContext) Pipeline() Pipeline {
	return c.channel.Pipeline()
}

func (c *DefaultHandlerContext) Write(msg interface{}) error {
	var (
		err error
		out = msg
	)
	return c.Pipeline().IterateOutbound(func(handler OutboundHandler) error {
		out, err = handler.Write(c, out)
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *DefaultHandlerContext) WriteAndFlush(msg interface{}) error {
	if err := c.Write(msg); err != nil {
		return err
	}
	if err := c.Flush(); err != nil {
		return err
	}
	return nil
}

func (c *DefaultHandlerContext) Flush() error {
	return c.channel.Flush()
}

func (c *DefaultHandlerContext) Close() error {
	return c.channel.Close()
}
