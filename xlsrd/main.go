package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	wkDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wkDir)
	xlsPath := filepath.Join(wkDir, "a.xls")
	fmt.Println(xlsPath)
	xlsFile, err := os.Open(xlsPath)
	if err != nil {
		log.Fatal(err)
	}
	head := make([]byte, 8)
	rC, err := xlsFile.Read(head)
	if err != nil {
		log.Fatal(err)
	}

	xlsHead := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}

	if bytes.Equal(xlsHead, head) {
		fmt.Println("有效 xls 文件头")
	} else {
		fmt.Printf("无效 xls 文件头 %d: %s \n", rC, hex.Dump(head))
	}
}
