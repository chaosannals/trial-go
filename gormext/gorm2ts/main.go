package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	hereDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(hereDir)
	modelsDir := filepath.Join(hereDir, "models")
	fmt.Println(modelsDir)
	modelsGlobStr := filepath.Join(modelsDir, "*.gen.go")
	fmt.Println(modelsGlobStr)
	files, err := filepath.Glob(modelsGlobStr)
	if err != nil {
		log.Fatalln(err)
	}
	for i, file := range files {
		fmt.Printf("gen: [%d] %s\n", i, file)
		if lexemes, err := readLexemes(file); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println(lexemes)
		}
	}
}
