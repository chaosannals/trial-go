package trial

type Conf struct {
	TcpHost  string
	TcpPort  uint16
	HttpHost string
	HttpPort uint16
	DbHost   string
	DbPort   uint16
	DbName   string
	DbUser   string
	DbPass   string
}

// Provider 必须返回 实例  可以有后续 error 返回从 provide 传出
func NewConf() (*Conf, error) {
	return &Conf{
		TcpHost:  "0.0.0.0",
		TcpPort:  44444,
		HttpHost: "0.0.0.0",
		HttpPort: 44440,
		DbHost:   "localhost",
		DbPort:   3306,
		DbName:   "exert",
		DbUser:   "root",
		DbPass:   "password",
	}, nil
}
