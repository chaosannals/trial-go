package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

var OPTS struct {
	InputPath string `short:"i" long:"input"`
	InputDir  string `short:"d" long:"dir"`
	LexOnly   bool   `short:"l" long:"lex"`
	ParseOnly bool   `short:"p" long:"parse"`
}

func makeByFile(inputPath string) []GoStruct {
	lexemes, err := readGoLexemes(inputPath)
	if err != nil {
		log.Fatalln(err)
	}
	if OPTS.LexOnly {
		for i, lexeme := range lexemes {
			fmt.Printf("%d %v\n", i, lexeme)
		}
		return []GoStruct{}
	}
	structs, err := parseGoStruct(lexemes)
	if err != nil {
		log.Fatalln(err)
	}
	if OPTS.ParseOnly {
		for i, s := range structs {
			fmt.Printf("%d %s\n", i, s.Name)
			for j, f := range s.Fields {
				fmt.Printf("%d\t %v\n", j, f)
			}
		}
		return []GoStruct{}
	}
	return structs
}

func main() {
	if _, err := flags.ParseArgs(&OPTS, os.Args); err != nil {
		log.Fatalln(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	n := 0
	inputPath := OPTS.InputPath
	if len(inputPath) > 0 {
		if !filepath.IsAbs(inputPath) {
			inputPath = filepath.Join(workDir, inputPath)
		}

		structs := makeByFile(inputPath)
		for i, s := range structs {
			fmt.Printf("%d ================== start\n", i)
			makeTs(&s)
			fmt.Printf("%d ================== end\n", i)
		}
		n += len(structs)
	}

	inputDir := OPTS.InputDir
	if len(inputDir) > 0 {
		if !filepath.IsAbs(inputDir) {
			inputDir = filepath.Join(workDir, inputDir)
		}
		pattern := filepath.Join(inputDir, "*.go")
		inputPaths, err := filepath.Glob(pattern)
		if err != nil {
			log.Fatalln(err)
		}
		for _, inputPath := range inputPaths {
			structs := makeByFile(inputPath)
			for i, s := range structs {
				fmt.Printf("%d ================== start\n", n+i)
				makeTs(&s)
				fmt.Printf("%d ================== end\n", n+i)
			}
			n += len(structs)
		}
	}
}
