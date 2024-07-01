package main

import (
	"fmt"
	"os"
	"unicode"
)

var keywords = map[string]bool{
	"package": true,
	"import":  true,
	"const":   true,
	"type":    true,
	"struct":  true,
	"func":    true,
	"return":  true,
}

type LexemeType = string

const (
	LEX_STRING  = "string"
	LEX_STRING2 = "string2"
	LEX_KEYWORD = "keyword"
	LEX_ID      = "id"
	LEX_PUNCT   = "punct"
	LEX_COMMENT = "comment"
	LEX_NL      = "next line"
)

type Lexeme struct {
	Type    LexemeType
	Content string
}

func readLexemes(modelSrcPath string) ([]Lexeme, error) {
	textBytes, err := os.ReadFile(modelSrcPath)
	if err != nil {
		return nil, err
	}
	text := []rune(string(textBytes))
	// fmt.Printf("rune: %v\n", text)

	lexer := &Lexer{
		index:  0,
		row:    1,
		column: 1,
		code:   text,
	}
	if err := lexer.lex(); err != nil {
		return nil, err
	}

	return lexer.result, nil
}

type Lexer struct {
	index  int
	row    int
	column int
	code   []rune
	result []Lexeme
}

func (lexer *Lexer) lex() error {
	for {
		start := lexer.index
		if start == len(lexer.code) {
			return nil
		}

		lexer.skipBlank()
		lexer.matchComment()
		lexer.matchPunct()
		lexer.matchString()
		lexer.matchIdOrKeyword()
		lexer.matchString2()
		end := lexer.index
		if start == end {
			fmt.Printf("词法不匹配： %d %d\n", lexer.row, lexer.column)
			break
		}
	}

	return nil
}

func (lexer *Lexer) getChar() rune {
	if len(lexer.code) <= lexer.index {
		return -1
	}
	c := lexer.code[lexer.index]
	return c
}

func (lexer *Lexer) peekChar(n int) rune {
	index := lexer.index + n
	if len(lexer.code) <= index {
		return -1
	}
	c := lexer.code[index]
	return c
}

func (lexer *Lexer) nextChar() rune {
	lexer.index++
	c := lexer.getChar()
	if c == rune('\n') {
		lexer.row++
		lexer.column = 1
		lexer.result = append(lexer.result, Lexeme{
			Type:    LEX_NL,
			Content: "\n",
		})
	} else {
		lexer.column++
	}
	return c
}

func (lexer *Lexer) matchPunct() error {
	c := lexer.getChar()
	if c == rune(')') {
		lexer.result = append(lexer.result, Lexeme{Content: ")", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune('(') {
		// fmt.Printf("===========cc: %v\n", c)
		lexer.result = append(lexer.result, Lexeme{Content: "(", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune('=') {
		lexer.result = append(lexer.result, Lexeme{Content: "=", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune('{') {
		lexer.result = append(lexer.result, Lexeme{Content: "{", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune('}') {
		lexer.result = append(lexer.result, Lexeme{Content: "}", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune('[') {
		lexer.result = append(lexer.result, Lexeme{Content: "[", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune(']') {
		lexer.result = append(lexer.result, Lexeme{Content: "]", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune('*') {
		lexer.result = append(lexer.result, Lexeme{Content: "*", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	} else if c == rune('.') {
		lexer.result = append(lexer.result, Lexeme{Content: ".", Type: LEX_PUNCT})
		lexer.nextChar()
		return nil
	}
	return nil
}

func (lexer *Lexer) matchComment() error {
	c1 := lexer.getChar()
	c2 := lexer.peekChar(1)
	if c1 != rune('/') || c2 != rune('/') {
		return nil
	}
	lexer.nextChar()
	c := lexer.nextChar()
	comment := []rune{}
	for c != rune('\n') {
		comment = append(comment, c)
		c = lexer.nextChar()
	}
	lexer.result = append(lexer.result, Lexeme{
		Type:    LEX_COMMENT,
		Content: string(comment),
	})
	return nil
}

func (lexer *Lexer) matchIdOrKeyword() error {
	c := lexer.getChar()
	if !unicode.IsLetter(c) {
		return nil
	}
	word := []rune{}
	for unicode.IsLetter(c) || unicode.IsDigit(c) {
		word = append(word, c)
		c = lexer.nextChar()
	}
	text := string(word)
	if _, ok := keywords[text]; ok {
		fmt.Printf("is keyword %s\n", text)
		lexer.result = append(lexer.result, Lexeme{
			Content: text,
		})
	} else {
		fmt.Printf("is id %s\n", text)
		lexer.result = append(lexer.result, Lexeme{
			Content: text,
		})
	}
	return nil
}

func (lexer *Lexer) matchString() error {
	c := lexer.getChar()
	if c != rune('"') {
		return nil
	}
	c = lexer.nextChar()
	word := []rune{}
	for c != rune('"') && c != -1 {
		word = append(word, c)
		c = lexer.nextChar()
	}
	lexer.nextChar()
	text := string(word)
	fmt.Printf("string: %s\n", text)
	lexer.result = append(lexer.result, Lexeme{
		Type:    LEX_STRING,
		Content: text,
	})
	return nil
}

func (lexer *Lexer) matchString2() error {
	c := lexer.getChar()
	if c != rune('`') {
		return nil
	}
	c = lexer.nextChar()
	word := []rune{}
	for c != rune('`') && c != -1 {
		word = append(word, c)
		c = lexer.nextChar()
	}
	lexer.nextChar()
	text := string(word)
	fmt.Printf("string2: %s\n", text)
	lexer.result = append(lexer.result, Lexeme{
		Type:    LEX_STRING2,
		Content: text,
	})
	return nil
}

func (lexer *Lexer) skipBlank() {
	for unicode.IsSpace(lexer.getChar()) {
		lexer.nextChar()
	}
}
