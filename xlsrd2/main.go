package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/chaosannals/xlsrd2/xlsrd2"
)

func main() {
	wkDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wkDir)
	xlsPath := filepath.Join(wkDir, "b.xls")
	xlsBook, err := xlsrd2.ReadXls(xlsPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("book struct size: %d\n", unsafe.Sizeof(*xlsBook))

	fmt.Println("======================================")
	wss, err := xlsBook.ListWookSheetNames()
	if err != nil {
		log.Fatal(err)
	}
	for i, ws := range wss {
		fmt.Printf("%d %s", i, ws.Name)
	}
}
