package main

import (
	"crypto/rc4"
	"fmt"
	"log"
)

// 这种算法带状态，所以每次都要重新创建 Cipher
// 密文再执行一次就解密
func DoRc4(key []byte, data []byte) ([]byte, error) {
	rc4Chiper, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 对称加密，输入和输出大小一致，没有补位。
	result := make([]byte, len(data))
	rc4Chiper.XORKeyStream(result, data)
	return data, nil
}

func main() {
	// RC4 算法的密钥长度是 256字节，不足会循环补充。
	// 例如：12345 的密钥会 1234512345123.... 直到补满
	// golang 的库会自动补满
	key := []byte("1234567890ABCDEF")
	text := "1234567890中文ABCDEF"
	data := []byte(text)
	fmt.Printf("text: %s\n", text)
	fmt.Printf("data: %v\n", data)

	result, err := DoRc4(key, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result: %v\n", result)

	resultBack, err := DoRc4(key, result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("resultBack: %v\n", resultBack)
	resultBackText := string(resultBack)
	fmt.Printf("resultBackText: %s\n", resultBackText)
}
