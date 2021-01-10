package channel

import (
	log "github.com/sirupsen/logrus"
	"reflect"
)

// 第一个入站处理器，直接调用下一个即可
type HeadInboundHandler struct {
	pipeline *DefaultPipeline
}

func NewHeadInboundHandler(pipeline *DefaultPipeline) *HeadInboundHandler {
	return &HeadInboundHandler{pipeline: pipeline}
}

func (p *HeadInboundHandler) OnError(c HandlerContext, err error) {
	c.FireOnError(err)
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

func (p *TailOutboundHandler) OnError(c HandlerContext, err error) {
	p.log.Warnln("unhandled error.", err)
}

func (p *TailOutboundHandler) Write(c HandlerContext, msg interface{}) error {
	buf, ok := msg.([]byte)
	if !ok {
		p.log.Warnln("call write with wrong msg type", reflect.ValueOf(msg).Type())
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

// 第一个出站处理器，占位用
type HeadOutboundHandler struct {
	pipeline *DefaultPipeline
}

func NewHeadOutboundHandler(pipeline *DefaultPipeline) *HeadOutboundHandler {
	return &HeadOutboundHandler{pipeline: pipeline}
}

// todo: 触发不到
func (p *HeadOutboundHandler) OnError(c HandlerContext, err error) {

}

func (p *HeadOutboundHandler) Write(c HandlerContext, msg interface{}) error {
	return c.Write(msg)
}

func (p *HeadOutboundHandler) Flush(c HandlerContext) error {
	return c.Flush()
}
