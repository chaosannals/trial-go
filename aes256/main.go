package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func HmacSha256Compute(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func Encrypt(key []byte, data interface{}) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	iv := make([]byte, block.BlockSize())
	plaintext := PKCS7Padding(bs, block.BlockSize())
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plaintext))
	encrypter.CryptBlocks(encrypted, plaintext)
	hmac := HmacSha256Compute(key, encrypted)
	result := bytes.NewBuffer(iv)
	result.Write(hmac)
	result.Write(encrypted)
	return base64.StdEncoding.EncodeToString(result.Bytes()), nil
}

func Decrypt(key []byte, text string, result interface{}) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return err
	}
	bbs := block.BlockSize()
	iv := data[:bbs]
	hmac := data[bbs : bbs+32]
	raw := data[bbs+32:]
	chmac := HmacSha256Compute(key, raw)
	if bytes.Compare(hmac, chmac) != 0 {
		return errors.New("hmac error.")
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(raw))
	decrypter.CryptBlocks(decrypted, raw)
	r := PKCS7UnPadding(decrypted)
	json.Unmarshal(r, result)
	return nil
}

type TestData struct {
	Name     string    `json:"name"`
	CreateAt time.Time `json:"createAt"`
}

func main() {
	k := []byte("PHQvHriVhxUJ2W7DiMip4dSPACxxA25F")
	a := TestData{
		Name:     "11111",
		CreateAt: time.Now(),
	}
	b, err := Encrypt(k, a)
	if err != nil {
		log.Fatalln(err)
	}
	r := &TestData{}
	err = Decrypt(k, b, r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r.Name)
	fmt.Println(r.CreateAt)
}
