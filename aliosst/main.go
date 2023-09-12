package main

import (
  "os"
  "fmt"
  "log"
  "encoding/json"
  "path/filepath"
  "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func main() {
  fmt.Println("OSS Go SDK Version: ", oss.Version)
  root, err := filepath.Abs(filepath.Dir(os.Args[0]))

  if err != nil {
    log.Fatalln(err)
  }

  confPath := filepath.Join(root, "alioss.conf.json")
}
