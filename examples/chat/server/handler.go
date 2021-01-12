package main

import (
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/examples/chat/common"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"sync"
)

var (
	userMap = make(map[string]string)
	locker  sync.RWMutex
)

type chatServerHandler struct {
	log *log.Entry
}

func (p chatServerHandler) OnError(c channel.HandlerContext, err error) {
	debug.PrintStack()
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}

func (p chatServerHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")
}

func (p chatServerHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p chatServerHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	p.log.Infof("receive %s", msg)
	packet := msg.(common.Packet)
	switch packet.Opcode() {
	case common.OpPing:
		p.handlePing(c)
	case common.OpRegisterReq:
		p.handleRegister(c, packet.(*common.RegisterReq))
	}
}

func (p chatServerHandler) handlePing(c channel.HandlerContext) {
	if err := c.WriteAndFlush(common.NewPongPacket()); err != nil {
		p.OnError(c, err)
	}
}

func (p chatServerHandler) handleRegister(c channel.HandlerContext, req *common.RegisterReq) {
	locker.RLock()
	_, exists := userMap[req.Username]
	locker.RUnlock()
	if exists {
		if err := c.WriteAndFlush(common.NewRegisterResult(1, "username is exists")); err != nil {
			p.OnError(c, err)
		}
		return
	}

	locker.Lock()
	userMap[req.Username] = req.Password
	locker.Unlock()
	if err := c.WriteAndFlush(common.NewRegisterResult(0, "ok")); err != nil {
		p.OnError(c, err)
	}
}

func newChatServerHandler() *chatServerHandler {
	return &chatServerHandler{log: log.WithField("component", "chatServerHandler")}
}
