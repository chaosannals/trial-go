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
	wss, err := xlsBook.ListWookSheetInfos()
	if err != nil {
		log.Fatal(err)
	}
	for i, ws := range wss {
		fmt.Printf(
			"[%d]\n name:'%s'\n LastColumnLetter:'%s'\n totalColumns: %d\n totalRows: %d\n",
			i,
			ws.Name,
			ws.LastColumnLetter,
			ws.TotalColumns,
			ws.TotalRows,
		)
	}
}
