package trial

import (
	"fmt"
	"net"
	"github.com/rs/zerolog"
)

type TcpServer struct {
	host string
	port uint16
	logger *zerolog.Logger
}

func NewTcpServer(conf *Conf, logger *zerolog.Logger) (*TcpServer, error) {
	return &TcpServer{
		host: conf.TcpHost,
		port: conf.TcpPort,
		logger: logger,
	}, nil
}

func (i *TcpServer) Serve() error {
	address := fmt.Sprintf("%s:%d", i.host, i.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	i.logger.Info().Msg("server start")
	for {
		conn, err := listener.Accept()
		if err != nil {
			i.logger.Error().Err(err).Msg("client accept failed.")
			continue
		}
		go func() {
			defer conn.Close()

		}()
	}
}
