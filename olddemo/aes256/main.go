package main

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/andreburgaud/crypt2go/ecb"
)

const (
	ENCODE_KEY = "1234567890123456"
)

type DataJsonInfo struct {
	String1 string `json:"string1"`
	String2 string `json:"string2"`
	Int1    uint32 `json:"int1"`
}

func PaddingZero(code []byte, c int) []byte {
	bc := len(code)
	mc := (bc % c)
	if mc == 0 {
		return code
	}
	pc := c - mc
	for i := 0; i < pc; i++ {
		code = append(code, 0)
	}
	return code
}

func (info *DataJsonInfo) Encode() (string, error) {
	raw, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	log.Println(string(raw))
	// 256 / 8 = 32 库按 key 长度使用 128 或 256 位算法
	block, err := aes.NewCipher(PaddingZero([]byte(ENCODE_KEY), 32))
	if err != nil {
		return "", err
	}
	mode := ecb.NewECBEncrypter(block)
	bs := PaddingZero(raw, mode.BlockSize())
	r := make([]byte, len(bs))
	mode.CryptBlocks(r, bs)
	return base64.StdEncoding.EncodeToString(r), nil
}

func main() {
	data := &DataJsonInfo{
		String1: "45678213456456",
		String2: "2022-10-29 12:00:12",
		Int1:    12354678,
	}
	log.Println(data.Encode())
}
