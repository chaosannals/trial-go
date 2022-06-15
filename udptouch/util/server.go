package util

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type TouchClientInfo struct {
	Ip        *net.IP
	LastTouch time.Time
}

type TouchServer struct {
	addr    *net.UDPAddr
	clients map[string]*TouchClientInfo
}

func NewTouchServer(port int) *TouchServer {
	return &TouchServer{
		addr: &net.UDPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: port,
		},
		clients: make(map[string]*TouchClientInfo),
	}
}

func (self *TouchServer) Serve() error {
	conn, err := net.ListenUDP("udp", self.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	var buffer [1024]byte
	var request TouchRequest
	for {
		n, client, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			fmt.Printf("read udp failed, err: %v", err)
			continue
		}
		if err := json.Unmarshal(buffer[:n], &request); err != nil {
			fmt.Printf("unmarshal request err: %v", err)
			continue
		}
		if _, ok := self.clients[request.Id]; ok {
			self.clients[request.Id].Ip = &client.IP
			self.clients[request.Id].LastTouch = time.Now()
		} else {
			self.clients[request.Id] = &TouchClientInfo{
				Ip:        &client.IP,
				LastTouch: time.Now(),
			}
		}

		response := &TouchResponse{
			YourIp: client.IP.String(),
		}
		if _, ok := self.clients[request.Target]; ok {
			response.TargetIp = self.clients[request.Target].Ip.String()
			response.TargetLastTouch = self.clients[request.Target].LastTouch
		}

		r, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("marshal response err: %v", err)
			continue
		}
		if n, err := conn.WriteToUDP(r, client); err != nil {
			fmt.Printf("response %d err: %v", n, err)
			continue
		}
	}
}
