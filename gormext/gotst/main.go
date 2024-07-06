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
	Separator bool   `short:"s" long:"separator"`
}

func makeFrom(inputPath string) []GoStruct {
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

func makeTsWith(inputPath string, n int) int {
	structs := makeFrom(inputPath)
	for i, s := range structs {
		if OPTS.Separator {
			fmt.Printf("%d ================== start\n", n+i)
		}
		makeTs(&s)
		if OPTS.Separator {
			fmt.Printf("%d ================== end\n", n+i)
		}
	}
	return len(structs) + n
}

func glob() ([]string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	result := []string{}
	inputPath := OPTS.InputPath
	if len(inputPath) > 0 {
		if !filepath.IsAbs(inputPath) {
			inputPath = filepath.Join(workDir, inputPath)
		}
		result = append(result, inputPath)
	}

	inputDir := OPTS.InputDir
	if len(inputDir) > 0 {
		if !filepath.IsAbs(inputDir) {
			inputDir = filepath.Join(workDir, inputDir)
		}
		pattern := filepath.Join(inputDir, "*.go")
		inputPaths, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		result = append(result, inputPaths...)
	}

	return result, nil
}

func main() {
	if _, err := flags.ParseArgs(&OPTS, os.Args); err != nil {
		log.Fatalln(err)
	}

	inputPaths, err := glob()
	if err != nil {
		log.Fatalln(err)
	}

	n := 0
	for _, inputPath := range inputPaths {
		n = makeTsWith(inputPath, n)
	}
}
