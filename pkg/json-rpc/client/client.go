package client

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Client struct {
	*rpc.Client
	conn net.Conn
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client: jsonrpc.NewClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
