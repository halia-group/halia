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
	"encoding/binary"
	"github.com/halia-group/halia/bootstrap"
	"github.com/halia-group/halia/channel"
	"github.com/halia-group/halia/examples/chat/common"
	"github.com/halia-group/halia/handler/codec"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	client := bootstrap.NewClient(&bootstrap.ClientOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().AddLast("frameDecoder", codec.NewLengthFieldBasedFrameDecoder(2, 4, binary.BigEndian))
			c.Pipeline().AddLast("packetDecoder", &common.PacketDecoder{})
			c.Pipeline().AddLast("packetEncoder", &common.PacketEncoder{})
			c.Pipeline().AddLast("handler", newChatClientHandler())

			return c
		},
	})

	log.WithField("component", "client").Fatal(client.Dial("tcp", "127.0.0.1:8080"))
}
