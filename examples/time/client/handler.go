package main

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"strings"
	"time"
)

type timeClientHandler struct {
	log *log.Entry
}

func NewTimeClientHandler() *timeClientHandler {
	return &timeClientHandler{log: log.WithField("component", "timeClientHandler")}
}

func (p timeClientHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}

func (p timeClientHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")

	p.log.Infof("pipeline in: %v", strings.Join(c.Pipeline().InboundNames(), "->"))
	p.log.Infof("pipeline out: %v", strings.Join(c.Pipeline().OutboundNames(), "->"))
}

func (p timeClientHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p timeClientHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	t := msg.(time.Time)
	p.log.Infoln("now is ", t.String())
}
