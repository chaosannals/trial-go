package trial

type Conf struct {
	TcpHost string
	TcpPort uint16
}

// Provider 必须返回 实例  可以有后续 error 返回从 provide 传出
func NewConf() (*Conf, error) {
	return &Conf{
		TcpHost: "127.0.0.1",
		TcpPort: 44444,
	}, nil
}