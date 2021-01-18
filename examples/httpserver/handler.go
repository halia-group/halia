package main

import (
	"errors"
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/handler/codec/http"
	log "github.com/sirupsen/logrus"
)

type handler struct{ log *log.Entry }

func (p handler) OnError(c channel.HandlerContext, err error) {
	log.WithError(err).Warnln("an error was caught")
}

func (p handler) ChannelActive(c channel.HandlerContext) {
	p.log.Infoln("connected")
}

func (p handler) ChannelInActive(c channel.HandlerContext) {
	p.log.Infoln("disconnected")
}

func (p handler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	r, ok := msg.(*http.Request)
	if !ok {
		p.OnError(c, errors.New("unknown request"))
		return
	}
	log.Infoln(r)
	resp := http.Response{
		Version: "HTTP/1.1",
		Status:  200,
		Reason:  "OK",
		Body:    []byte("Hello World"),
	}
	if err := c.WriteAndFlush(&resp); err != nil {
		p.OnError(c, err)
	}
}

func newHandler() *handler {
	return &handler{log: log.WithField("component", "handler")}
}
