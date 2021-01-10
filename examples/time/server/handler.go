package main

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"reflect"
	"strings"
	"time"
)

type timeServerHandler struct {
	log    *log.Entry
	ticker *time.Ticker
}

func newTimeServerHandler() *timeServerHandler {
	return &timeServerHandler{log: log.WithField("component", "timeServerHandler")}
}

func (p *timeServerHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("error caught %v %v \n", err, reflect.ValueOf(err).Type())
	p.ticker.Stop()
}

func (p *timeServerHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")
	p.log.Infof("pipeline in: %v", strings.Join(c.Pipeline().InboundNames(), "->"))
	p.log.Infof("pipeline out: %v", strings.Join(c.Pipeline().OutboundNames(), "->"))

	p.ticker = time.NewTicker(time.Second)
	for range p.ticker.C {
		if err := c.WriteAndFlush(time.Now()); err != nil {
			p.OnError(c, err)
		}
	}
}

func (p timeServerHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p timeServerHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {

}
