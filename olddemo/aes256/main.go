package main

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
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

	log.Printf("encode: raw len: %d", len(raw))

	// bs := PaddingZero(raw, mode.BlockSize())

	// 这个库还提供了一些补齐方法。
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	bs, err := padder.Pad(raw)
	if err != nil {
		return "", err
	}

	log.Printf("encode: bs len: %d %% %d \n", len(bs), mode.BlockSize())

	r := make([]byte, len(bs))
	mode.CryptBlocks(r, bs)

	log.Printf("encode: len: %d \n", len(r))
	return base64.StdEncoding.EncodeToString(r), nil
}

func Decode(encodeText string) (*DataJsonInfo, error) {
	raw, err := base64.StdEncoding.DecodeString(encodeText)
	if err != nil {
		return nil, err
	}

	log.Printf("decode: len: %d \n", len(raw))

	block, err := aes.NewCipher(PaddingZero([]byte(ENCODE_KEY), 32))
	if err != nil {
		return nil, err
	}

	mode := ecb.NewECBDecrypter(block)
	bs := make([]byte, len(raw))
	mode.CryptBlocks(bs, raw)

	log.Printf("decode: bs len: %d %% %d \n", len(bs), mode.BlockSize())

	padder := padding.NewPkcs7Padding(mode.BlockSize())
	r, err := padder.Unpad(bs)

	if err != nil {
		return nil, err
	}

	// r := bs[:len(bs)-4]

	log.Printf("decode: r len: %d", len(r))

	var di DataJsonInfo
	if err := json.Unmarshal(r, &di); err != nil {
		return nil, err
	}
	return &di, nil
}

func main() {
	data := &DataJsonInfo{
		String1: "45678213456456",
		String2: "2022-10-29 12:00:12",
		Int1:    12354678,
	}
	ed,err := data.Encode()
	if err !=nil {
		log.Printf("error1: %v \n", err)
	}
	log.Println(ed)

	r, err := Decode(ed)
	if err !=nil {
		log.Printf("error2: %v \n", err)
	}
	log.Println(r)
}
