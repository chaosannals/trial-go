package stress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type HttpStressWorker struct {
	config *StressConfig
	client *http.Client
}

func NewHttpStressWorker(config *StressConfig) *HttpStressWorker {
	return &HttpStressWorker{
		config: config,
		client: &http.Client{},
	}
}

func (self *HttpStressWorker) RequestGet() {
	target := fmt.Sprintf(
		"%s://%s:%d/%s",
		self.config.Scheme,
		self.config.Host,
		self.config.Port,
		self.config.Path,
	)
	r, err := http.Get(target)
	if err != nil {
		fmt.Printf("get err: %v", err)
	}
	defer r.Body.Close()
}

func (self *HttpStressWorker) RequestPost() ([]byte, error) {
	target := fmt.Sprintf(
		"%s://%s:%d/%s",
		self.config.Scheme,
		self.config.Host,
		self.config.Port,
		self.config.Path,
	)
	body, err := json.Marshal(self.config.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("body %s\n", string(body))
	req, err := http.NewRequest(
		"POST",
		target,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("res body: %s", string(data))
	return data, nil
}

func (self *HttpStressWorker) Request() {
	var wg sync.WaitGroup
	wg.Add(self.config.Times)
	for i := 0; i < self.config.Times; i += 1 {
		j := i
		go func() {
			// fmt.Printf("invoker: %d", j)
			start := time.Now()
			switch self.config.Method {
			case "GET":
				self.RequestGet()
				break
			case "POST":
				_, err := self.RequestPost()
				if err != nil {
					fmt.Printf("no: %d err %v\n", j, err)
				}
				break
			}
			elapsed := time.Since(start)
			fmt.Printf("elapsed: %f s\n", elapsed.Seconds())
			wg.Done()
		}()
	}
	wg.Wait()
}
