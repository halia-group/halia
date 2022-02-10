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
