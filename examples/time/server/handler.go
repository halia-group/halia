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
	log "github.com/sirupsen/logrus"
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
