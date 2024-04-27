package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chaosannals/xlsrd3/xlsrd3"
)

func main() {
	wkDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wkDir)
	xlsPath := filepath.Join(wkDir, "b.xls")
	fmt.Printf("load xls file at: %s \n", xlsPath)
	xlsBook, err := xlsrd3.LoadFormFile(xlsPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v \n", xlsBook)
}
