package main

import (
	"github.com/halia-group/halia/channel"
	log "github.com/sirupsen/logrus"
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
	p.log.Infoln("pipeline", c.Pipeline().Names())
}

func (p *EchoServerHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p *EchoServerHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	c.FireChannelRead(msg)
	//data, ok := msg.([]byte)
	//if !ok {
	//	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("unknown msg type: %+v", msg)
	//	return
	//}
	//str := string(data)
	//p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("receive ", str)
	//if err := c.WriteAndFlush("server:" + str + "\r\n"); err != nil {
	//	p.log.WithField("peer", c.Channel().RemoteAddr()).WithError(err).Warnln("write error")
	//}
}
