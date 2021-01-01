package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
	"io"
	"net"
)

type Server struct {
	listener net.Listener
	options  *ServerOptions
	log      *log.Entry
}

func NewServer(options *ServerOptions) *Server {
	return &Server{options: options, log: log.WithField("component", "server")}
}

func (server *Server) Listen(network, addr string) (err error) {
	server.listener, err = net.Listen(network, addr)
	if err != nil {
		return
	}
	var (
		conn net.Conn
	)
	server.log.WithField("network", network).WithField("addr", addr).Infoln("started")
	for {
		conn, err = server.listener.Accept()
		if err != nil {
			return
		}
		go server.onConnect(conn)
	}
}

func (server *Server) onConnect(conn net.Conn) {
	var (
		c = server.options.HandlerContextFactory(server.options.ChannelFactory(conn))
	)
	defer server.onDisconnect(c)

	c.Pipeline().IterateInbound(func(handler channel.InboundHandler) error {
		handler.ChannelActive(c)
		return nil
	})

	var (
		buf = make([]byte, 4096)
	)
	for {
		n, err := c.Channel().Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			server.onReadError(c, err)
		}
		server.onRead(c, buf[:n])
	}
}

func (server *Server) onDisconnect(c channel.HandlerContext) {
	c.Pipeline().IterateInbound(func(handler channel.InboundHandler) error {
		handler.ChannelInActive(c)
		return nil
	})
}

func (server *Server) onReadError(c channel.HandlerContext, err error) {
	c.Pipeline().IterateInbound(func(handler channel.InboundHandler) error {
		handler.ErrorCaught(c, err)
		return nil
	})
}

func (server *Server) onRead(c channel.HandlerContext, data []byte) {
	var (
		msg interface{} = data
		err error
	)
	c.Pipeline().IterateInbound(func(handler channel.InboundHandler) error {
		msg, err = handler.ChannelRead(c, msg)
		if err != nil {
			return err
		}
		return nil
	})
}
