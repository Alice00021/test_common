package client

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func NewClient(address string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewClient(conn), nil
}
