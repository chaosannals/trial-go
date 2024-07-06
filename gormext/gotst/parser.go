package main

import (
	"fmt"
	"time"
)

type GoField struct {
	Name     string
	Type     []string
	Tag      string
	Comment  string
	Comment2 string
}

type GoStruct struct {
	Name   string
	Fields []GoField
}

func parseGoStruct(lexemes []GoLexeme) ([]GoStruct, error) {
	parser := &GoStructParser{
		index:   0,
		lexemes: lexemes,
		current: nil,
		results: []GoStruct{},
	}

	cancel := goTimeout(time.Second*4, func() {
		fmt.Printf("parse at: %d %v\n", parser.index, parser.peekLexeme(0))
	})
	defer cancel()

	err := parser.parse()
	return parser.results, err
}

type GoStructParser struct {
	index   int
	lexemes []GoLexeme
	current *GoStruct
	results []GoStruct
}

func (parser *GoStructParser) parse() error {
	for parser.peekLexeme(0) != nil {
		parser.skipIgnore()
		if err := parser.matchStruct(); err != nil {
			return err
		}
	}
	return nil
}

func (parser *GoStructParser) matchStruct() error {
	if parser.peekLexeme(0) == nil {
		return nil
	}

	if !parser.matchKeyword(KW_TYPE, 0) {
		return fmt.Errorf("不是预期的 type 关键字：%v", parser.peekLexeme(0))
	}
	name := parser.nextLexeme()
	if name == nil || name.Type != LEX_ID {
		return fmt.Errorf("不是预期的标识符：%v", name)
	}
	parser.popLexemes(1)
	parser.current = &GoStruct{
		Name:   name.Content,
		Fields: []GoField{},
	}
	if !parser.matchKeyword(KW_STRUCT, 0) {
		return fmt.Errorf("不是预期的 struct 关键字：%v", parser.peekLexeme(0))
	}
	parser.popLexemes(1)

	if !parser.matchPunt("{", 0) {
		return fmt.Errorf("不是预期的操作符 { ：%v", parser.peekLexeme(0))
	}
	parser.popLexemes(1)

	for {
		if nextOk, err := parser.matchField(); err != nil {
			return err
		} else {
			if !nextOk {
				break
			}
		}
	}

	if !parser.matchPunt("}", 0) {
		return fmt.Errorf("不是预期的操作符 } ：%v", parser.peekLexeme(0))
	}
	parser.popLexemes(1)

	parser.results = append(parser.results, *parser.current)
	return nil
}

func (parser *GoStructParser) matchField() (bool, error) {
	name := parser.peekLexeme(0)
	if name.Type != LEX_ID {
		return false, nil
	}
	parser.popLexemes(1)

	field := &GoField{
		Name: name.Content,
	}

	if err := parser.matchFieldType(field); err != nil {
		return false, err
	}

	string2 := parser.peekLexeme(0)
	if string2.Type == LEX_STRING2 {
		field.Tag = string2.Content
		parser.nextLexeme()
	}

	comment := parser.peekLexeme(0)
	if comment.Type == LEX_COMMENT {
		field.Comment = comment.Content
		parser.nextLexeme()
	}

	comment2 := parser.peekLexeme(0)
	if comment2.Type == LEX_COMMENT2 {
		field.Comment2 = comment2.Content
		parser.nextLexeme()
	}

	// fmt.Printf("field: %v\n", field)
	parser.current.Fields = append(parser.current.Fields, *field)

	return true, nil
}

func (parser *GoStructParser) matchFieldType(field *GoField) error {
	if parser.matchPunt("[", 0) {
		parser.nextLexeme()
		if !parser.matchPunt("]", 0) {
			return fmt.Errorf("语法错误 %v", parser.peekLexeme(0))
		}
		parser.nextLexeme()
		field.Type = append(field.Type, "[]")
		return parser.matchFieldType(field)
	} else if parser.matchPunt("*", 0) {
		parser.nextLexeme()
		field.Type = append(field.Type, "*")
		return parser.matchFieldType(field)
	} else if parser.matchPunt(".", 0) {
		parser.nextLexeme()
		field.Type = append(field.Type, ".")
		return parser.matchFieldType(field)
	} else {
		id := parser.peekLexeme(0)
		if id.Type == LEX_ID {
			parser.nextLexeme()
			field.Type = append(field.Type, id.Content)
			return parser.matchFieldType(field)
		}
		// else {
		// 	return fmt.Errorf("类型语法错误 %v", id.Content)
		// }
	}
	return nil
}

func (parser *GoStructParser) matchPunt(punt string, n int) bool {
	lexeme := parser.peekLexeme(n)
	return (lexeme != nil) &&
		(lexeme.Type == LEX_PUNCT) &&
		(lexeme.Content == punt)
}

// func (parser *GoStructParser) matchId(n int) bool {
// 	lexeme := parser.peekLexeme(n)
// 	return (lexeme != nil) &&
// 		(lexeme.Type == LEX_ID)
// }

func (parser *GoStructParser) matchKeyword(kw GoLangKw, n int) bool {
	lexeme := parser.peekLexeme(n)
	return (lexeme != nil) &&
		(lexeme.Type == LEX_KEYWORD) &&
		(lexeme.Content == string(kw))
}

func (parser *GoStructParser) skipIgnore() {
	for !parser.matchKeyword(KW_TYPE, 0) {
		if parser.nextLexeme() == nil {
			break
		}
	}
}

func (parser *GoStructParser) peekLexeme(n int) *GoLexeme {
	index := parser.index + n
	if len(parser.lexemes) <= index {
		return nil
	}
	return &parser.lexemes[index]
}

func (parser *GoStructParser) nextLexeme() *GoLexeme {
	parser.index++
	if len(parser.lexemes) <= parser.index {
		return nil
	}
	return &parser.lexemes[parser.index]
}

func (parser *GoStructParser) popLexemes(n int) {
	for i := 0; i < n; i++ {
		parser.nextLexeme()
	}
}
