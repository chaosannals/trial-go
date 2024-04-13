package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Utf8ToGbk(text string) (string, error) {
	r := bytes.NewReader([]byte(text))
	decoder := transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	result, err := io.ReadAll(decoder)
	fmt.Printf("result: %d\n", len(result))
	return string(result), err
}

func main() {
	fmt.Println("start")
	a := "ROOT ENTRY 中文"
	b, err := Utf8ToGbk(a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("a: %d b: %d\n", len(a), len(b))
	fmt.Printf("a: %s b: %s\n", a, b)
	fmt.Printf("len: %d\n", len("SUMMARYINFORMATION"))
}
