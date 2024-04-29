package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chaosannals/xlsrd4/xlsrd4"
)

func main() {
	wkDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wkDir)
	xlsPath := filepath.Join(wkDir, "b.xls")
	fmt.Printf("load xls file at: %s \n", xlsPath)
	// xlsrd4.ReadXlsFile(xlsPath)
	sheets, err := xlsrd4.ListXlsSheetInfo(xlsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, sheet := range sheets {
		fmt.Printf("Sheet: %s \n", sheet.Name)
	}
}
