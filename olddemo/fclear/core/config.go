package core

type FClearAliyunOssConfig struct {
	Tag string
	Host string
	Key string
	Secret string
	Bucket string
	Endpoint string
	Prefix string
	IsStorageIA bool
	IsPrivate bool
}

type FClearConfigTask struct {
	Dir string
	Exts []string
	ValidTime string
	Alioss string
}

type FClearConfig struct {
	Mode string
	Tasks []FClearConfigTask
	Alioss []FClearAliyunOssConfig
}
