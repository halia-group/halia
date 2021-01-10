package main

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"strings"
)

type EchoServerHandler struct {
	log *log.Entry
}

func NewEchoServerHandler() *EchoServerHandler {
	return &EchoServerHandler{
		log: log.WithField("component", "EchoServerHandler"),
	}
}

func (p *EchoServerHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}

func (p *EchoServerHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")

	p.log.Infof("pipeline in: %v", strings.Join(c.Pipeline().InboundNames(), "->"))
	p.log.Infof("pipeline out: %v", strings.Join(c.Pipeline().OutboundNames(), "->"))
}

func (p *EchoServerHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p *EchoServerHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	str, ok := msg.(string)
	if !ok {
		p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("unknown msg type: %+v", msg)
		return
	}
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("receive ", str)
	if err := c.Write("server:" + str + "\r\n"); err != nil {
		p.log.WithField("peer", c.Channel().RemoteAddr()).WithError(err).Warnln("write error")
	}
}
