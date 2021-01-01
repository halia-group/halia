package main

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
)

type EchoServerHandler struct {
	channel.InboundHandlerAdapter
	log *log.Entry
}

func NewEchoServerHandler() *EchoServerHandler {
	return &EchoServerHandler{log: log.WithField("component", "EchoServerHandler")}
}

func (p EchoServerHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Debugln("connected")
}

func (p EchoServerHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Debugln("disconnected")
}

func (p EchoServerHandler) ChannelRead(c channel.HandlerContext, msg interface{}) (out interface{}, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Debugf("read: %+v", msg)
	return
}
