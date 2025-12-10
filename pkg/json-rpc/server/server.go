package server

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	address  string
	server   *rpc.Server
	listener net.Listener
}

func NewServer(address string) *Server {
	s := &Server{
		address: address,
		server:  rpc.NewServer(),
	}

	return s
}

func (s *Server) Register(service interface{}) error {
	return s.server.Register(service)
}

func (s *Server) StartTCP(errChan chan<- error) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	s.listener = listener

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if s.listener == nil {
					return
				}
				errChan <- err
				return
			}
			go s.server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
	return nil
}

func (s *Server) Shutdown() error {
	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		return err
	}
	return nil
}
