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

func (server *Server) Listen(network, addr string) error {
	var err error
	server.listener, err = net.Listen(network, addr)
	if err != nil {
		return err
	}
	server.log.WithField("network", network).WithField("addr", addr).Infoln("started")
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			return err
		}
		go server.onConnect(conn)
	}
}

// todo: 如何在读取失败的时候返回
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
