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

package bootstrap

import (
	"net"

	"github.com/halia-group/halia/channel"
)

type ServerOptions struct {
	ChannelFactory func(conn net.Conn) channel.Channel
}

type Server struct {
	listener net.Listener
	options  *ServerOptions
}

func NewServer(options *ServerOptions) *Server {
	return &Server{options: options}
}

func (server *Server) Listen(network, addr string) error {
	var err error
	server.listener, err = net.Listen(network, addr)
	if err != nil {
		return err
	}
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			return err
		}
		go server.onConnect(conn)
	}
}

func (server *Server) onConnect(conn net.Conn) {
	c := server.options.ChannelFactory(conn)
	defer server.onDisconnect(c)

	c.Pipeline().FireChannelActive()
	// 数据包读取由入站handler进行轮询读取
	c.Pipeline().FireChannelRead(nil)
}

// 断开连接
func (server *Server) onDisconnect(c channel.Channel) {
	c.Pipeline().FireChannelInActive()
}
