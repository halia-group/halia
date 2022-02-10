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
