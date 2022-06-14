package stress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
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
	rb := PreloadMake(self.config.Body)
	// fmt.Printf("rb: %v\n", rb)

	body, err := json.Marshal(rb)
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

	if len(self.config.AcceptEncoding) > 0 {
		req.Header.Set("Accept-Encoding", self.config.AcceptEncoding)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// data, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return nil, err
	// }

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		switch self.config.AcceptEncoding {
		case "br":
			//fmt.Printf("start br: %v\n", string(data))
			//n := len(data)
			// b := bytes.NewReader(data)
			// reader := brotli.NewReader(b)
			reader := brotli.NewReader(res.Body)
			return ioutil.ReadAll(reader)
			// d, err := ioutil.ReadAll(reader)
			// if err != nil {
			// 	return nil, err
			// }
			// fmt.Printf("br: %d => %d\n", n, len(body))
			// return d, nil
		case "gzip":
			reader, err := gzip.NewReader(res.Body)
			if err != nil {
				return nil, err
			}
			return ioutil.ReadAll(reader)
		case "flate":
			reader := flate.NewReader(res.Body)
			return ioutil.ReadAll(reader)
		default:
			return ioutil.ReadAll(res.Body)
		}
		// return data, nil
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("res %d body: %s", res.StatusCode, string(data))
	return data, fmt.Errorf("http response error: %d\n", res.StatusCode)
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
				// else {
				// 	fmt.Printf("data: %s", string(d))
				// }
				break
			}
			elapsed := time.Since(start)
			fmt.Printf("elapsed: %f s\n", elapsed.Seconds())
			wg.Done()
		}()
	}
	wg.Wait()
}
