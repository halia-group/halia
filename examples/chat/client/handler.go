package main

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"halia/examples/chat/common"
)

type chatClientHandler struct {
	log *log.Entry
}

func (p chatClientHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}

func (p chatClientHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")
	if err := c.WriteAndFlush(common.NewPingPacket()); err != nil {
		p.OnError(c, err)
	}
}

func (p chatClientHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

func (p chatClientHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	p.log.Infof("receive %s", msg)

	packet := msg.(common.Packet)
	switch packet.Opcode() {
	case common.OpPong:
		p.handlePong(c)
	case common.OpRegisterResult:
		p.handleRegisterResult(c, packet.(*common.RegisterResult))
	}
}

func (p chatClientHandler) handlePong(c channel.HandlerContext) {
	packet := common.NewRegisterReq("xialei", "111111")
	if err := c.WriteAndFlush(packet); err != nil {
		p.OnError(c, err)
	}
}

func (p chatClientHandler) handleRegisterResult(c channel.HandlerContext, result *common.RegisterResult) {
}

func newChatClientHandler() *chatClientHandler {
	return &chatClientHandler{log: log.WithField("component", "chatClientHandler")}
}
