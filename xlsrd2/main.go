package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chaosannals/xlsrd2/xlsrd2"
)

func main() {
	wkDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wkDir)
	xlsPath := filepath.Join(wkDir, "b.xls")
	err = xlsrd2.ReadXls(xlsPath)
	if err != nil {
		log.Fatal(err)
	}
}
