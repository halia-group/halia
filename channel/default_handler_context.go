/*
 *
 *  * MIT License
 *  *
 *  * Copyright (c) [2021] [xialeistudio]
 *  *
 *  * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  * of this software and associated documentation files (the "Software"), to deal
 *  * in the Software without restriction, including without limitation the rights
 *  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  * copies of the Software, and to permit persons to whom the Software is
 *  * furnished to do so, subject to the following conditions:
 *  *
 *  * The above copyright notice and this permission notice shall be included in all
 *  * copies or substantial portions of the Software.
 *  *
 *  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  * SOFTWARE.
 *
 */

package channel

type defaultHandlerContext struct {
	next     *defaultHandlerContext
	prev     *defaultHandlerContext
	pipeline *defaultPipeline
	name     string
	handler  Handler
}

func (c *defaultHandlerContext) findNextInbound() *defaultHandlerContext {
	ctx := c
	for ctx.next != nil {
		ctx = ctx.next
		if _, ok := ctx.handler.(InboundHandler); ok {
			break
		}
	}
	return ctx
}
func (c *defaultHandlerContext) findNextOutbound() *defaultHandlerContext {
	ctx := c
	for ctx.prev != nil {
		ctx = ctx.prev
		if _, ok := ctx.handler.(OutboundHandler); ok {
			break
		}
	}
	return ctx
}

func (c *defaultHandlerContext) FireChannelActive() {
	var (
		next    = c.findNextInbound()
		handler = next.handler.(InboundHandler)
	)
	handler.ChannelActive(next)
}

func (c *defaultHandlerContext) FireChannelInActive() {
	var (
		next    = c.findNextInbound()
		handler = next.handler.(InboundHandler)
	)
	handler.ChannelInActive(next)
}

func (c *defaultHandlerContext) FireChannelRead(msg interface{}) {
	var (
		next    = c.findNextInbound()
		handler = next.handler.(InboundHandler)
	)
	handler.ChannelRead(next, msg)
}

func (c *defaultHandlerContext) FireOnError(err error) {
	var (
		next    = c.findNextInbound()
		handler = next.handler.(InboundHandler)
	)
	handler.OnError(next, err)
}

func (c *defaultHandlerContext) Write(msg interface{}) error {
	var (
		prev    = c.findNextOutbound()
		handler = prev.handler.(OutboundHandler)
	)
	return handler.Write(prev, msg)
}

func (c *defaultHandlerContext) Flush() error {
	var (
		prev    = c.findNextOutbound()
		handler = prev.handler.(OutboundHandler)
	)
	return handler.Flush(prev)
}

func (c *defaultHandlerContext) WriteAndFlush(msg interface{}) error {
	var (
		prev    = c.findNextOutbound()
		handler = prev.handler.(OutboundHandler)
	)
	if err := handler.Write(prev, msg); err != nil {
		return err
	}
	return handler.Flush(prev)
}

func (c *defaultHandlerContext) Channel() Channel {
	return c.pipeline.channel
}

func (c *defaultHandlerContext) Name() string {
	return c.name
}

func (c *defaultHandlerContext) Handler() Handler {
	return c.handler
}

func (c *defaultHandlerContext) Pipeline() Pipeline {
	return c.pipeline
}

func (c *defaultHandlerContext) OnError(err error) {
	var (
		next    = c.next
		handler = next.Handler().(InboundHandler)
	)
	handler.OnError(next, err)
}
