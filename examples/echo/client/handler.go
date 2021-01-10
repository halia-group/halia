package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"strings"
	"time"
)

type EchoClientHandler struct {
	log *log.Entry
}

func NewEchoClientHandler() *EchoClientHandler {
	return &EchoClientHandler{
		log: log.WithField("component", "EchoClientHandler"),
	}
}

func (p *EchoClientHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}

func (p *EchoClientHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")

	if err := c.WriteAndFlush("Hello World\r\n"); err != nil {
		p.log.WithError(err).Warnln("write error")
	}
	p.log.Infof("pipeline in: %v", strings.Join(c.Pipeline().InboundNames(), "->"))
	p.log.Infof("pipeline out: %v", strings.Join(c.Pipeline().OutboundNames(), "->"))
}

func (p *EchoClientHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p *EchoClientHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	str, ok := msg.(string)
	if !ok {
		p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("unknown msg type: %+v", msg)
		return
	}
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("receive ", str)
	time.AfterFunc(time.Second, func() {
		if err := c.WriteAndFlush(fmt.Sprintf("client say:%s\r\n", time.Now().String())); err != nil {
			p.log.WithError(err).Warnln("write error")
		}
	})
}
