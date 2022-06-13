package stress

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type StressConfig struct {
	Scheme string
	Host   string
	Port   int
	Method string
	Path   string
	Times  int
	Body   interface{}
}

func LoadConfig(path string) (*StressConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var config StressConfig
	content, err := ioutil.ReadAll(f)
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
