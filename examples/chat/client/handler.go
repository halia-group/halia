/*
 *
 *  * MIT License
 *  *
 *  * Copyright (c) [2021] [xialeistudio]
 *  *
 *  * Permission is hereby granted, free of charge, to any person obtaining a copy
 *  * of this software and associated documentation files (the "Software"), to deal
 *  * in the Software without restriction, including without limitation the rights
 *  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  * copies of the Software, and to permit persons to whom the Software is
 *  * furnished to do so, subject to the following conditions:
 *  *
 *  * The above copyright notice and this permission notice shall be included in all
 *  * copies or substantial portions of the Software.
 *  *
 *  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  * SOFTWARE.
 *
 */

package main

import (
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/examples/chat/common"
	log "github.com/sirupsen/logrus"
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
