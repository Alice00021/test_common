package client

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type TCPClient struct {
	*rpc.Client
	conn net.Conn
}

func NewTCPClient(address string) (*TCPClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPClient{
		Client: jsonrpc.NewClient(conn),
		conn:   conn,
	}, nil
}

func (c *TCPClient) Close() error {
	return c.conn.Close()
}

func NewClient(address string) (*TCPClient, error) {
	return NewTCPClient(address)
}

func NewHTTPClient(host, path string) (*rpc.Client, error) {
	return rpc.DialHTTPPath("tcp", host, path)
}
