package server

import (
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	address string
}

func NewServer(address string) *Server {
	return &Server{address: address}
}

func (s *Server) Register(service interface{}) error {
	return rpc.Register(service)
}

func (s *Server) Start() error {
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonrpc.ServeConn(conn)
	})

	return http.ListenAndServe(s.address, nil)
}
