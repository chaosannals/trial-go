package util

import (
	"encoding/json"
	"fmt"
	"net"
)

type TouchClient struct {
	server     *net.UDPAddr
	identifier string
	target     string
}

func NewTouchClient(server string, identifier string, target string) (*TouchClient, error) {
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return nil, err
	}
	return &TouchClient{
		server:     addr,
		identifier: identifier,
		target:     target,
	}, nil
}

func (self *TouchClient) Start() error {
	conn, err := net.DialUDP("udp", nil, self.server)
	if err != nil {
		return err
	}
	defer conn.Close()
	response, err := self.touchServe(conn)
	if err != nil {
		return err
	}

	fmt.Printf("response: %v \n", response)

	return nil
}

func (self *TouchClient) touchServe(conn *net.UDPConn) (*TouchResponse, error) {
	request := &TouchRequest{
		Id:     self.identifier,
		Target: self.target,
	}
	req, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(req); err != nil {
		return nil, err
	}
	var buffer [1024]byte
	n, err := conn.Read(buffer[:])
	if err != nil {
		return nil, err
	}
	var response TouchResponse
	if err := json.Unmarshal(buffer[:n], &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (self *TouchClient) scanTargetPort(conn *net.UDPConn) {
	
}