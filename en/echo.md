# Halia

Halia is a component-based network application framework for rapid development of maintainable high-performance servers
and clients.

Reference netty to implementation. [Netty](https://netty.io/)

[toc]

## Get Started

This example will demonstrate how to develop a time echo server.

Client sends a local time string to server every 1 second, and the server echoes the data.

### Public Code

#### encoder.go

The String-Encoder, converts `string` to `[]byte` and then transfers to the next outbound processor.

```go
package main

import (
	"halia/channel"
)

type StringToByteEncoder struct{}

// The encoder doesn't process error ，transfers to the next processor (the business processor) for processing
func (e *StringToByteEncoder) OnError(c channel.HandlerContext, err error) {
	c.FireOnError(err)
}

func (e *StringToByteEncoder) Write(c channel.HandlerContext, msg interface{}) error {
	// if msg is string, we do convert and transfer to next processor.
	if str, ok := msg.(string); ok { 
		return c.Write([]byte(str))
	}
	return c.Write(msg)
}

func (e *StringToByteEncoder) Flush(c channel.HandlerContext) error {
	return c.Flush()
}
```

### Client Code

#### handler.go

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

// This callback is called when an error occurred.
func (p *EchoClientHandler) OnError(c channel.HandlerContext, err error) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Warnln("error caught", err)
}

// This callback is called when connection has been established.
func (p *EchoClientHandler) ChannelActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("connected")

	if err := c.WriteAndFlush("Hello World\r\n"); err != nil {
		p.log.WithError(err).Warnln("write error")
	}
	p.log.Infof("pipeline in: %v", c.Pipeline().Names())
}

// This callback is called when connection has been closed.
func (p *EchoClientHandler) ChannelInActive(c channel.HandlerContext) {
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("disconnected")
}

// This callback is called when we read a complete message.
func (p *EchoClientHandler) ChannelRead(c channel.HandlerContext, msg interface{}) {
	data, ok := msg.([]byte)
	if !ok {
		p.log.WithField("peer", c.Channel().RemoteAddr()).Warnf("unknown msg type: %+v", msg)
		return
	}
	str := string(data)
	p.log.WithField("peer", c.Channel().RemoteAddr()).Infoln("receive ", str)
	// after 1 second, we send data to server
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
		// Wrap the raw `net.Conn` to `Channel`, generally, you can use the `DefaultChannel`.
		ChannelFactory: func(conn net.Conn) channel.Channel {
			c := channel.NewDefaultChannel(conn)
			// add decoder
			c.Pipeline().AddLast("decoder", codec.NewLineBasedFrameDecoder())
			// add encoder
			c.Pipeline().AddLast("encoder", &StringToByteEncoder{})
			// add business processor
			c.Pipeline().AddLast("handler", NewEchoClientHandler())
			return c
		},
	})
	// connect server
	log.WithField("component", "client").Fatal(client.Dial("tcp", "127.0.0.1:8080"))
}
```

### Server Code

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

### Run

Run the server first, then run the client

**Server Output**

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

**Client Output**

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