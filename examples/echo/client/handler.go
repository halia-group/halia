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
