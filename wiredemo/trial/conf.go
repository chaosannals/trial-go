package trial

type Conf struct {
	TcpHost string
	TcpPort uint16
}

// Provider 必须以 (实例，error) 或 (实例, clearup, error) 返回
// clearup 提供回收方法
func NewConf() (*Conf, error) {
	return &Conf{
		TcpHost: "127.0.0.1",
		TcpPort: 44444,
	}, nil
}