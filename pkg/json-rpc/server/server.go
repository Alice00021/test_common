package server

import (
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	address string
	server  *rpc.Server
}

func NewServer(address string) *Server {
	s := &Server{
		address: address,
		server:  rpc.NewServer(),
	}

	http.HandleFunc("/rpc", s.handleRPC)
	http.HandleFunc("/health", s.handleHealth)

	return s
}

func (s *Server) Register(service interface{}) error {
	return s.server.Register(service)
}

func (s *Server) StartTCP() error {
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

func (s *Server) StartHTTP() error {
	return http.ListenAndServe(s.address, nil)
}

func (s *Server) Start() error {
	return s.StartHTTP()
}

func (s *Server) handleRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	conn := &httpConn{
		r:      r.Body,
		w:      w,
		closed: false,
	}

	s.server.ServeCodec(jsonrpc.NewServerCodec(conn))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

type httpConn struct {
	r      io.ReadCloser
	w      io.Writer
	closed bool
}

func (c *httpConn) Read(p []byte) (n int, err error) {
	if c.closed {
		return 0, io.EOF
	}
	return c.r.Read(p)
}

func (c *httpConn) Write(p []byte) (n int, err error) {
	if c.closed {
		return 0, io.EOF
	}
	return c.w.Write(p)
}

func (c *httpConn) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.r.Close()
}
