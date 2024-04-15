package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func Utf8ToGbk(text string) (string, error) {
	r := bytes.NewReader([]byte(text))
	decoder := transform.NewReader(r, simplifiedchinese.GBK.NewEncoder())
	result, err := io.ReadAll(decoder)
	fmt.Printf("result: %d\n", len(result))
	return string(result), err
}

func GbkToUtf8(text string) (string, error) {
	r := bytes.NewReader([]byte(text))
	decoder := transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	result, err := io.ReadAll(decoder)
	fmt.Printf("result: %d\n", len(result))
	return string(result), err
}

func Utf8ToUtf16Le(text string) (string, error) {
	r := bytes.NewReader([]byte(text))
	decoder := transform.NewReader(r, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder())
	result, err := io.ReadAll(decoder)
	fmt.Printf("result: %d\n", len(result))
	return string(result), err
}

func Utf16LeToUtf8(text string) (string, error) {
	result, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), []byte(text))
	return string(result), err
}

func main() {
	fmt.Println("start")
	utf8Text := "ROOT ENTRY 中文"
	gbkText, err := Utf8ToGbk(utf8Text)
	if err != nil {
		log.Fatal(err)
	}
	gbkTextBack, err := GbkToUtf8(gbkText)
	if err != nil {
		log.Fatal(err)
	}
	utf16LeText, err := Utf8ToUtf16Le(utf8Text)
	if err != nil {
		log.Fatal(err)
	}
	utf16LeTextBack, err := Utf16LeToUtf8(utf16LeText)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("utf8Text: %d %s\n", len(utf8Text), utf8Text)
	fmt.Printf("utf16LeText: %d %s\n", len(utf16LeText), utf16LeText)
	fmt.Printf("utf16LeTextBack: %d %s\n", len(utf16LeTextBack), utf16LeTextBack)
	fmt.Printf("gbkText: %d %s\n", len(gbkText), gbkText)
	fmt.Printf("gbkTextBack: %d %s\n", len(gbkTextBack), gbkTextBack)
	fmt.Printf("len: %d\n", len("SUMMARYINFORMATION"))
}
