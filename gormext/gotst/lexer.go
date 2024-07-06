package main

import (
	"fmt"
	"os"
	"time"
	"unicode"
)

type GoLangKw string
type GoLangLex = string

const (
	KW_PACKAGE GoLangKw = "package"
	KW_IMPORT  GoLangKw = "import"
	KW_CONST   GoLangKw = "const"
	KW_TYPE    GoLangKw = "type"
	KW_STRUCT  GoLangKw = "struct"
	KW_FUNC    GoLangKw = "func"
	KW_RETURN  GoLangKw = "return"
)

const (
	LEX_STRING   GoLangLex = "string"
	LEX_STRING2  GoLangLex = "string2"
	LEX_KEYWORD  GoLangLex = "keyword"
	LEX_ID       GoLangLex = "id"
	LEX_PUNCT    GoLangLex = "punct"
	LEX_COMMENT  GoLangLex = "comment"
	LEX_COMMENT2 GoLangLex = "comment2"
)

const (
	RUNE_EOF rune = -1
	RUNE_NL  rune = '\n'
)

var GOL_KW_SET = map[GoLangKw]bool{
	KW_PACKAGE: true,
	KW_IMPORT:  true,
	KW_CONST:   true,
	KW_TYPE:    true,
	KW_STRUCT:  true,
	KW_FUNC:    true,
	KW_RETURN:  true,
}

var GOL_PUNT_SET = map[rune][][]rune{
	'+': {{'+'}, {'='}, {}},
	'-': {{'-'}, {'='}, {}},
	'*': {{'='}, {}},
	'/': {{'='}, {}},
	'%': {{'='}, {}},
	'!': {{'='}, {}},
	'=': {{'='}, {}},
	'<': {{'<', '='}, {'='}, {'<'}, {'-'}, {}},
	'>': {{'>', '='}, {'='}, {'>'}, {}},
	'&': {{'&', '='}, {'&'}, {'='}, {}},
	'|': {{'|', '='}, {'|'}, {'='}, {}},
	'^': {{'='}, {}},
	'(': {{}},
	')': {{}},
	'{': {{}},
	'}': {{}},
	'[': {{}},
	']': {{}},
	'.': {{}},
}

type GoLexeme struct {
	Type    GoLangLex
	Content string
}

func readGoLexemes(srcPath string) ([]GoLexeme, error) {
	srcBytes, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	src := []rune(string(srcBytes))

	lexer := &GoLexer{
		index:  0,
		row:    1,
		column: 1,
		code:   src,
		err:    nil,
	}

	cancel := goTimeout(time.Second*4, func() {
		fmt.Printf("[%s]at: %d %d\n", srcPath, lexer.row, lexer.column)
	})
	defer cancel()

	err = lexer.lex()
	return lexer.result, err
}

type GoLexer struct {
	index  int
	row    int
	column int
	code   []rune
	result []GoLexeme
	err    error
}

func (lexer *GoLexer) lex() error {
	for lexer.peekChar(0) != RUNE_EOF {
		if err := lexer.
			skipBlank().
			matchComment().
			matchComment2().
			matchPunct().
			matchIdOrKeyword().
			matchString().
			matchString2().
			err; err != nil {
			return err
		}
	}
	return lexer.err
}

func (lexer *GoLexer) matchPunctNext(chars []rune) bool {
	for i, nextChar := range chars {
		peekChar := lexer.peekChar(i)
		if nextChar != peekChar {
			return false
		}
	}
	return true
}

func (lexer *GoLexer) matchPunct() *GoLexer {
	if lexer.err != nil {
		return lexer
	}

	c := lexer.peekChar(0)
	if nextCharsSet, ok := GOL_PUNT_SET[c]; ok {
		lexer.nextChar()
		for _, nextChars := range nextCharsSet {
			if lexer.matchPunctNext(nextChars) {
				lexer.popChars(len(nextChars))
				lexer.result = append(lexer.result, GoLexeme{
					Content: string(append([]rune{c}, nextChars...)),
					Type:    LEX_PUNCT,
				})
				return lexer
			}
		}
		lexer.err = fmt.Errorf("非法 GOLANG 操作符 %v", c)
	}
	return lexer
}

func (lexer *GoLexer) matchComment() *GoLexer {
	if lexer.err != nil {
		return lexer
	}

	c1 := lexer.peekChar(0)
	c2 := lexer.peekChar(1)
	if c1 != '/' || c2 != '/' {
		return lexer
	}
	lexer.popChars(2)

	c := lexer.nextChar()
	comment := []rune{}
	for c != RUNE_NL {
		comment = append(comment, c)
		c = lexer.nextChar()
	}

	lexer.result = append(lexer.result, GoLexeme{
		Type:    LEX_COMMENT,
		Content: string(comment),
	})
	return lexer
}

func (lexer *GoLexer) matchComment2() *GoLexer {
	if lexer.err != nil {
		return lexer
	}

	c1 := lexer.peekChar(0)
	c2 := lexer.peekChar(1)
	if c1 != '/' || c2 != '*' {
		return lexer
	}
	lexer.popChars(2)

	c := lexer.nextChar()
	comment2 := []rune{}
	for c != RUNE_EOF {
		if c == '*' && lexer.peekChar(1) == '/' {
			lexer.popChars(2)
			break
		}
		comment2 = append(comment2, c)
		c = lexer.nextChar()
	}

	lexer.result = append(lexer.result, GoLexeme{
		Type:    LEX_COMMENT2,
		Content: string(comment2),
	})
	return lexer
}

func (lexer *GoLexer) matchIdOrKeyword() *GoLexer {
	c := lexer.peekChar(0)
	if !unicode.IsLetter(c) {
		return lexer
	}

	word := []rune{}
	for unicode.IsLetter(c) || unicode.IsDigit(c) {
		word = append(word, c)
		c = lexer.nextChar()
	}

	text := string(word)
	if _, ok := GOL_KW_SET[GoLangKw(text)]; ok {
		// fmt.Printf("is keyword %s\n", text)
		lexer.result = append(lexer.result, GoLexeme{
			Type:    LEX_KEYWORD,
			Content: text,
		})
	} else {
		// fmt.Printf("is id %s\n", text)
		lexer.result = append(lexer.result, GoLexeme{
			Type:    LEX_ID,
			Content: text,
		})
	}
	return lexer
}

func (lexer *GoLexer) matchString() *GoLexer {
	c := lexer.peekChar(0)
	if c != '"' {
		return lexer
	}

	c = lexer.nextChar()
	word := []rune{}
	for c != '"' && c != RUNE_EOF {
		if c == RUNE_NL {
			lexer.err = fmt.Errorf("字符串不可换行[%d %d]", lexer.row, lexer.column)
			return lexer
		}
		if c == '\\' {
			switch lexer.peekChar(0) {
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 't':
				c = '\t'
			case '"':
				c = '"'
			case '\'':
				c = '\''
			case '\\':
				c = '\\'
			default:
				lexer.err = fmt.Errorf("非法转义符号： %v", c)
				return lexer
			}
		}
		word = append(word, c)
		c = lexer.nextChar()
	}
	lexer.popChars(1)

	text := string(word)
	// fmt.Printf("string: %s\n", text)
	lexer.result = append(lexer.result, GoLexeme{
		Type:    LEX_STRING,
		Content: text,
	})
	return lexer
}

func (lexer *GoLexer) matchString2() *GoLexer {
	c := lexer.peekChar(0)
	if c != '`' {
		return lexer
	}

	c = lexer.nextChar()
	word := []rune{}
	for c != '`' && c != RUNE_EOF {
		word = append(word, c)
		c = lexer.nextChar()
	}
	lexer.popChars(1)

	text := string(word)
	// fmt.Printf("string2: %s\n", text)
	lexer.result = append(lexer.result, GoLexeme{
		Type:    LEX_STRING2,
		Content: text,
	})
	return lexer
}

func (lexer *GoLexer) skipBlank() *GoLexer {
	for unicode.IsSpace(lexer.peekChar(0)) {
		lexer.nextChar()
	}
	return lexer
}

func (lexer *GoLexer) peekChar(n int) rune {
	index := lexer.index + n
	if len(lexer.code) <= index {
		return RUNE_EOF
	}
	c := lexer.code[index]
	return c
}

func (lexer *GoLexer) nextChar() rune {
	lexer.index++
	c := lexer.peekChar(0)
	if c == RUNE_NL {
		lexer.row++
		lexer.column = 1
		// fmt.Printf("换行 %d", lexer.row)
	} else {
		lexer.column++
	}
	return c
}

func (lexer *GoLexer) popChars(n int) {
	for i := 0; i < n; i++ {
		lexer.nextChar()
	}
}
