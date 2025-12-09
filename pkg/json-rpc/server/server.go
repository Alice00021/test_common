package server

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	address string
	server  *rpc.Server
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		server:  rpc.NewServer(),
	}
}

func (s *Server) Register(service interface{}) error {
	return s.server.Register(service)
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
