package main

import (
	"fmt"
	"github.com/halia-group/halia/channel"
	log "github.com/sirupsen/logrus"
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

	p.log.Infoln("pipeline", c.Pipeline().Names())
	if err := c.WriteAndFlush("Hello World\r\n"); err != nil {
		p.log.WithError(err).Warnln("write error")
	}
}

func (p *EchoClientHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p *EchoClientHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	data, ok := msg.([]byte)
	if !ok {
		p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("unknown msg type: %+v", msg)
		return
	}
	str := string(data)
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("receive ", str)
	time.AfterFunc(time.Second, func() {
		if err := c.WriteAndFlush(fmt.Sprintf("client say:%s\r\n", time.Now().String())); err != nil {
			p.log.WithError(err).Warnln("write error")
		}
	})
}
