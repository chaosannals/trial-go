package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

var OPTS struct {
	InputPath string `short:"i" long:"input" required:"true"`
	LexOnly   bool   `short:"l" long:"lex"`
	ParseOnly bool   `short:"p" long:"parse"`
}

func main() {
	if _, err := flags.ParseArgs(&OPTS, os.Args); err != nil {
		log.Fatalln(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	inputPath := OPTS.InputPath
	if !filepath.IsAbs(inputPath) {
		inputPath = filepath.Join(workDir, inputPath)
	}

	lexemes, err := readGoLexemes(inputPath)
	if err != nil {
		log.Fatalln(err)
	}
	if OPTS.LexOnly {
		for i, lexeme := range lexemes {
			fmt.Printf("%d %v\n", i, lexeme)
		}
	}
	structs, err := parseGoStruct(lexemes)
	if err != nil {
		log.Fatalln(err)
	}
	for i, s := range structs {
		fmt.Printf("%d %s\n", i, s.Name)
		for j, f := range s.Fields {
			fmt.Printf("%d\t %v\n", j, f)
		}
	}
}
