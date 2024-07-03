package main

import (
	"fmt"
)

type FieldInfo struct {
	Name    string
	Type    []string
	Comment string
}

type TypeInfo struct {
	Name   string
	Fields []FieldInfo
}

type Parser struct {
	index   int
	lexemes []Lexeme
	current *TypeInfo
	Results []TypeInfo
}

func (parser *Parser) ParseType() error {
	for parser.peekLexeme(0) != nil {
		parser.matchTypeAndSkipIgnore()
		name := parser.nextLexeme()
		if name == nil {
			break
		}
		fmt.Printf("=====================start %s\n", name.Content)
		parser.current = &TypeInfo{
			Name:   name.Content,
			Fields: []FieldInfo{},
		}
		fmt.Printf("type: %s\n", name.Content)
		if !parser.matchKeyword(KW_STRUCT, 1) {
			return fmt.Errorf("struct 不匹配")
		} else {
			fmt.Println("struct")
		}
		parser.nextLexeme()
		if !parser.matchPunt("{", 1) {
			return fmt.Errorf("{ 不匹配")
		} else {
			fmt.Println("{")
		}
		parser.nextLexeme()
		parser.nextLexeme()
		for parser.matchField() {
			// fmt.Println("matchField")
		}
		if !parser.matchPunt("}", 0) {
			return fmt.Errorf("} 不匹配")
		} else {
			fmt.Println("}")
		}
		parser.nextLexeme()
		parser.Results = append(parser.Results, *parser.current)
		fmt.Printf("=====================end %v\n", name.Content)
	}
	return nil
}

func (parser *Parser) peekLexeme(n int) *Lexeme {
	index := parser.index + n
	if len(parser.lexemes) <= index {
		return nil
	}
	return &parser.lexemes[index]
}

func (parser *Parser) nextLexeme() *Lexeme {
	parser.index++
	if len(parser.lexemes) <= parser.index {
		return nil
	}
	return &parser.lexemes[parser.index]
}

func (parser *Parser) matchKeyword(kw KEYWORD, n int) bool {
	lexeme := parser.peekLexeme(n)
	// fmt.Println(lexeme)
	v := (lexeme != nil) && (lexeme.Type == LEX_KEYWORD) && (lexeme.Content == string(kw))
	// fmt.Printf("=============vvv %s == %s = %v\n", kw, lexeme.Content, v)
	return v
}

func (parser *Parser) matchPunt(punt string, n int) bool {
	lexeme := parser.peekLexeme(n)
	return (lexeme != nil) && (lexeme.Type == LEX_PUNCT) && (lexeme.Content == punt)
}

func (parser *Parser) matchField() bool {
	name := parser.peekLexeme(0)
	fmt.Printf("start matchField %s  %s\n", name.Content, name.Type)
	if name.Type != LEX_ID {
		return false
	}

	parser.nextLexeme()
	field := &FieldInfo{
		Name: name.Content,
	}

	if !parser.matchFieldType(field) {
		return false
	}

	string2 := parser.peekLexeme(0)
	if string2.Type == LEX_STRING2 {
		field.Comment = string2.Content
		parser.nextLexeme()
	}

	comment := parser.peekLexeme(0)
	if comment.Type == LEX_COMMENT {
		parser.nextLexeme()
	}

	fmt.Printf("field: %v\n", field)
	parser.current.Fields = append(parser.current.Fields, *field)

	return true
}

func (parser *Parser) matchFieldType(field *FieldInfo) bool {
	type1 := parser.peekLexeme(0)
	if type1.Type == LEX_PUNCT && type1.Content == "[" {
		parser.nextLexeme()
		type2 := parser.peekLexeme(0)
		if type2.Type != LEX_PUNCT && type2.Content != "]" {
			fmt.Printf("语法错误 %v", type2.Content)
			return false
		}
		parser.nextLexeme()
		field.Type = append(field.Type, "[]")
		return parser.matchFieldType(field)
	} else if type1.Type == LEX_PUNCT && type1.Content == "*" {
		parser.nextLexeme()
		field.Type = append(field.Type, "*")
		return parser.matchFieldType(field)
	} else if type1.Type == LEX_PUNCT && type1.Content == "." {
		parser.nextLexeme()
		field.Type = append(field.Type, ".")
		return parser.matchFieldType(field)
	} else if type1.Type == LEX_ID {
		parser.nextLexeme()
		field.Type = append(field.Type, type1.Content)
		return parser.matchFieldType(field)
	} else {
		// fmt.Printf("类型匹配不上。%v", type1.Content)
	}
	return true
}

func (parser *Parser) matchTypeAndSkipIgnore() {
	for !parser.matchKeyword(KW_TYPE, 0) {
		if parser.nextLexeme() == nil {
			break
		}
	}
}
