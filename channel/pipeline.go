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

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Pipeline interface {
	InboundInvoker
	OutboundInvoker

	AddFirst(name string, handler Handler)
	AddLast(name string, handler Handler)
	Channel() Channel
	Names() []string
}

type defaultPipeline struct {
	head    *defaultHandlerContext
	tail    *defaultHandlerContext
	channel Channel
}

func NewDefaultPipeline(channel Channel) *defaultPipeline {
	pipeline := &defaultPipeline{
		channel: channel,
	}
	headCtx := &defaultHandlerContext{
		pipeline: pipeline,
		name:     "head",
		handler:  &headHandler{},
	}
	tailCtx := &defaultHandlerContext{
		pipeline: pipeline,
		name:     "tail",
		handler:  &tailHandler{},
	}
	headCtx.next = tailCtx
	tailCtx.prev = headCtx

	pipeline.head = headCtx
	pipeline.tail = tailCtx
	return pipeline
}

func (p *defaultPipeline) FireChannelActive() {
	p.head.FireChannelActive()
}

func (p *defaultPipeline) FireChannelInActive() {
	p.head.FireChannelInActive()
}

func (p *defaultPipeline) FireChannelRead(msg interface{}) {
	p.head.FireChannelRead(msg)
}

func (p *defaultPipeline) FireOnError(err error) {
	p.head.FireOnError(err)
}

func (p *defaultPipeline) Write(msg interface{}) error {
	return p.tail.Write(msg)
}

func (p *defaultPipeline) Flush() error {
	return p.tail.Flush()
}

func (p *defaultPipeline) WriteAndFlush(msg interface{}) error {
	return p.tail.WriteAndFlush(msg)
}

func (p *defaultPipeline) AddFirst(name string, handler Handler) {
	newCtx := &defaultHandlerContext{
		pipeline: p,
		name:     name,
		handler:  handler,
	}
	currentNext := p.head.next
	// connect currentNext and newCtx
	newCtx.next = currentNext
	currentNext.prev = newCtx
	// connect head and newCtx
	p.head.next = newCtx
	newCtx.prev = p.head
}

func (p *defaultPipeline) AddLast(name string, handler Handler) {
	newCtx := &defaultHandlerContext{
		pipeline: p,
		name:     name,
		handler:  handler,
	}
	currentPrev := p.tail.prev
	// connect currentPrev and newCtx
	newCtx.prev = currentPrev
	currentPrev.next = newCtx
	// connect tail and newCtx
	p.tail.prev = newCtx
	newCtx.next = p.tail
}

func (p *defaultPipeline) Names() []string {
	result := make([]string, 0)
	cursor := p.head
	for cursor != nil {
		result = append(result, cursor.name)
		cursor = cursor.next
	}
	return result
}

func (p *defaultPipeline) Channel() Channel {
	return p.channel
}

type headHandler struct{}

func (p headHandler) Write(c HandlerContext, msg interface{}) error {
	switch data := msg.(type) {
	case []byte:
		_, err := c.Channel().Write(data)
		return err
	default:
		return errors.New("write wrong msg type to head")
	}
}

func (p headHandler) Flush(c HandlerContext) error {
	return c.Channel().Flush()
}

func (p headHandler) OnError(c HandlerContext, err error) {
	c.FireOnError(err)
}

func (p headHandler) ChannelActive(c HandlerContext) {
	c.FireChannelActive()
}

func (p headHandler) ChannelInActive(c HandlerContext) {
	c.FireChannelInActive()
}

func (p headHandler) ChannelRead(c HandlerContext, msg interface{}) {
	c.FireChannelRead(msg)
}

type tailHandler struct{}

func (p tailHandler) ChannelActive(c HandlerContext) {

}

func (p tailHandler) ChannelInActive(c HandlerContext) {
}

func (p tailHandler) ChannelRead(c HandlerContext, _ interface{}) {
	log.WithField("component", "TailHandler").Debug("unhandled message that reached at the tail of the pipeline")
}

func (p tailHandler) OnError(c HandlerContext, err error) {
	log.WithField("component", "TailHandler").Debugf("unhandled error(%v) that reached at the tail of the pipeline", err)
}

func (p tailHandler) Write(c HandlerContext, msg interface{}) error {
	log.WithField("component", "TailHandler").Debug("unhandled write that reached at the tail of the pipeline")
	return nil
}

func (p tailHandler) Flush(c HandlerContext) error {
	log.WithField("component", "TailHandler").Debug("unhandled flush that reached at the tail of the pipeline")
	return nil
}
