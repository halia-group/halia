package bootstrap

import (
	log "github.com/sirupsen/logrus"
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

func (server *Server) onConnect(conn net.Conn) {

}
