# Halia

Halia是一个基于组件化的网络应用框架，用于快速开发可维护的高性能协议服务器和客户端。

## 快速开始

本示例将演示如何开发一个时间回显服务器。

客户端每隔1秒发送时间字符串给服务器，服务器回显该数据。

### 公用代码

#### encoder.go

字符串编码器，将字符串转换为`[]byte`传输到下一个出站处理器

```go
package main

import (
	"halia/channel"
)

type StringToByteEncoder struct{}
// 编码器不处理错误，交由下一个处理器(也就是业务处理器)处理
func (e *StringToByteEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (e *StringToByteEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	if str, ok := msg.(string); ok { // string才转换
		return c.Write([]byte(str))
	}
	return c.Write(msg)
}

func (e *StringToByteEncoder) Flush(c channel.HandlerContext) error {
	return c.Flush()
}
```

### 客户端代码

#### handler.go

客户端业务处理代码。

```go
package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"strings"
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
// 发送错误回调
func (p *EchoClientHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}
// 连接已建立
func (p *EchoClientHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")

	if err := c.WriteAndFlush("Hello World\r\n"); err != nil {
		p.log.WithError(err).Warnln("write error")
	}
	p.log.Infof("pipeline in: %v", c.Pipeline().Names())
}

// 连接已断开
func (p *EchoClientHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

// 读取到完整的消息回调
func (p *EchoClientHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	data, ok := msg.([]byte)
	if !ok {
		p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("unknown msg type: %+v", msg)
		return
	}
	str := string(data)
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("receive ", str)
    // 1秒后发送数据给服务器
	time.AfterFunc(time.Second, func() {
		if err := c.WriteAndFlush(fmt.Sprintf("client say:%s\r\n", time.Now().String())); err != nil {
			p.log.WithError(err).Warnln("write error")
		}
	})
}
```

#### main.go

```go
package main

import (
	log "github.com/sirupsen/logrus"
	"halia/bootstrap"
	"halia/channel"
	"halia/handler/codec"
	"net"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	client := bootstrap.NewClient(&bootstrap.ClientOptions{
        // 将原始net.Conn包装为Channel实现，一般情况下用DefaultChannel即可
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
            // 添加解码器，换行符分割报文解码器
			c.Pipeline().AddLast("decoder", codec.NewLineBasedFrameDecoder())
            // 添加编码器
			c.Pipeline().AddLast("encoder", &StringToByteEncoder{})
			// 添加业务处理器
			c.Pipeline().AddLast("handler", NewEchoClientHandler())
			return c
		},
	})
	// 连接服务器
	log.WithField("component", "client").Fatal(client.Dial("tcp", "127.0.0.1:8080"))
}
```

### 服务端代码

#### handler.go

```go
package main

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"strings"
)

type EchoServerHandler struct {
	log *log.Entry
}

func NewEchoServerHandler() *EchoServerHandler {
	return &EchoServerHandler{
		log: log.WithField("component", "EchoServerHandler"),
	}
}

func (p *EchoServerHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}

func (p *EchoServerHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")
}

func (p *EchoServerHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")

	p.log.Infof("pipeline in: %v", strings.Join(c.Pipeline().InboundNames(), "->"))
	p.log.Infof("pipeline out: %v", strings.Join(c.Pipeline().OutboundNames(), "->"))
}

func (p *EchoServerHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	data, ok := msg.([]byte)
	if !ok {
		p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("unknown msg type: %+v", msg)
		return
	}
	str := string(data)
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("receive ", str)
	if err := c.Write("server:" + str + "\r\n"); err != nil {
		p.log.WithField("peer", c.Channel().RemoteAddr()).WithError(err).Warnln("write error")
	}
}
```

#### main.go

```go
package main

import (
	log "github.com/sirupsen/logrus"
	"halia/bootstrap"
	"halia/channel"
	"halia/handler/codec"
	"net"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	s := bootstrap.NewServer(&bootstrap.ServerOptions{
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			c.Pipeline().AddLast("decoder", codec.NewLineBasedFrameDecoder())
			c.Pipeline().AddLast("encoder", &StringToByteEncoder{})
			c.Pipeline().AddLast("handler", NewEchoServerHandler())
			return c
		},
	})

	log.WithField("component", "server").Fatal(s.Listen("tcp", "0.0.0.0:8080"))
}
```

### 运行

先运行服务端，再运行客户端。

服务端输出

```text
time="2021-01-12T11:30:13+08:00" level=info msg=started addr="0.0.0.0:8080" component=server network=tcp pid=7584
time="2021-01-12T11:30:13+08:00" level=info msg=initialized component=channelId machineId=a0c5895a25a3 pid=7584
time="2021-01-12T11:30:18+08:00" level=info msg=connected component=EchoServerHandler peer="127.0.0.1:57641"
time="2021-01-12T11:30:18+08:00" level=info msg="pipeline in: InHeadContext->decoder->handler" component=EchoServerHandler
time="2021-01-12T11:30:18+08:00" level=info msg="pipeline out: OutHeadContext->encoder->OutTailContext" component=EchoServerHandler
time="2021-01-12T11:30:18+08:00" level=info msg="receive  Hello World" component=EchoServerHandler peer="127.0.0.1:57641"
time="2021-01-12T11:30:19+08:00" level=info msg="receive  client say:2021-01-12 11:30:19.5192868 +0800 CST m=+1.046443501" component=EchoServerHandler peer="127.0.0.1:57641"
time="2021-01-12T11:30:20+08:00" level=info msg="receive  client say:2021-01-12 11:30:20.5193884 +0800 CST m=+2.046545101" component=EchoServerHandler peer="127.0.0.1:57641"
time="2021-01-12T11:30:21+08:00" level=info msg="receive  client say:2021-01-12 11:30:21.5345887 +0800 CST m=+3.061745401" component=EchoServerHandler peer="127.0.0.1:57641"
time="2021-01-12T11:30:22+08:00" level=info msg="receive  client say:2021-01-12 11:30:22.5459978 +0800 CST m=+4.073154501" component=EchoServerHandler peer="127.0.0.1:57641"
```

客户端输出

```
time="2021-01-12T11:30:18+08:00" level=info msg=connected component=EchoClientHandler peer="127.0.0.1:8080"
time="2021-01-12T11:30:18+08:00" level=info msg="pipeline in: InHeadContext->decoder->handler" component=EchoClientHandler
time="2021-01-12T11:30:18+08:00" level=info msg="pipeline out: OutHeadContext->encoder->OutTailContext" component=EchoClientHandler
time="2021-01-12T11:30:18+08:00" level=info msg="receive  server:Hello World" component=EchoClientHandler peer="127.0.0.1:8080"
time="2021-01-12T11:30:18+08:00" level=info msg=initialized component=channelId machineId=a0c5895a25a3 pid=960
time="2021-01-12T11:30:19+08:00" level=info msg="receive  server:client say:2021-01-12 11:30:19.5192868 +0800 CST m=+1.046443501" component=EchoClientHandler peer="127.0.0.1:8080"
time="2021-01-12T11:30:20+08:00" level=info msg="receive  server:client say:2021-01-12 11:30:20.5193884 +0800 CST m=+2.046545101" component=EchoClientHandler peer="127.0.0.1:8080"
time="2021-01-12T11:30:21+08:00" level=info msg="receive  server:client say:2021-01-12 11:30:21.5345887 +0800 CST m=+3.061745401" component=EchoClientHandler peer="127.0.0.1:8080"
time="2021-01-12T11:30:22+08:00" level=info msg="receive  server:client say:2021-01-12 11:30:22.5459978 +0800 CST m=+4.073154501" component=EchoClientHandler peer="127.0.0.1:8080"
```





``