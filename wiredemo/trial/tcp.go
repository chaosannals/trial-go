package trial

import (
	"fmt"
	"log"
	"net"
)

type TcpServer struct {
	host string
	port uint16
}

func NewTcpServer(conf *Conf) (*TcpServer, error) {
	return &TcpServer{
		host: conf.TcpHost,
		port: conf.TcpPort,
	}, nil
}

func (i *TcpServer) Serve() error {
	address := fmt.Sprintf("%s:%d", i.host, i.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("client error: %v", err)
			continue
		}
		go func() {
			defer conn.Close()

		}()
	}
}
