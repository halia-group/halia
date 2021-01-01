package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"halia/channel"
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

	// todo: 暂时不读取
	select {}
}

func (server *Server) onDisconnect(c channel.HandlerContext) {
	c.Pipeline().IterateInbound(func(handler channel.InboundHandler) error {
		handler.ChannelInActive(c)
		return nil
	})
}
