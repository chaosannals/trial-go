package stress

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type PreloadParser struct {
	charIndex  int
	charRow    int
	charColumn int
	charCount  int
	chars      []rune
}

func (self *PreloadParser) Parse() (interface{}, error) {
	return self.matchExpression()
}

func (self *PreloadParser) matchExpression() (interface{}, error) {
	return self.matchFunctionCall()
}

func (self *PreloadParser) matchFunctionCall() (interface{}, error) {
	var p interface{}
	id := self.matchIdentifier()
	if err := self.match("("); err != nil {
		return nil, err
	}

	if c := self.peekChar(); c != ')' {
		if pv, err := self.matchExpression(); err != nil {
			return nil, err
		} else {
			p = pv
		}
	} else {
		p = "nil"
	}

	if err := self.match(")"); err != nil {
		return nil, err
	}

	t := reflect.TypeOf(self)
	m, b := t.MethodByName(id)
	if b {
		mp := make([]reflect.Value, 2)
		mp[0] = reflect.ValueOf(self)
		mp[1] = reflect.ValueOf(p)
		mr := m.Func.Call(mp)
		// fmt.Printf("call %s(%v) -> %v \n", id, p, mr[0])
		return mr[0].Interface(), nil
	} else {
		return nil, fmt.Errorf("method %s not found", id)
	}
}

func (self *PreloadParser) matchIdentifier() string {
	start := self.charIndex
	for i := start; i < self.charCount; i++ {
		if !unicode.IsLetter(self.chars[i]) {
			self.charIndex = i
			return string(self.chars[start:i])
		}
	}
	return string(self.chars[start:])
}

func (self *PreloadParser) match(token string) error {
	start := self.charIndex
	end := start + len(token)
	m := string(self.chars[start:end])
	if m == token {
		self.charIndex = end
		return nil
	}
	return fmt.Errorf("%s not match %s", m, token)
}

func (self *PreloadParser) popChar() rune {
	c := self.chars[self.charIndex]
	self.charIndex++
	return c
}

func (self *PreloadParser) peekChar() rune {
	return self.chars[self.charIndex]
}

func (self *PreloadParser) RandomObject(_ interface{}) interface{} {
	result := make(map[string]interface{})
	result["test"] = "111111"
	result["bbb"] = 1232
	return result
}

func (self *PreloadParser) Json(input interface{}) interface{} {
	// fmt.Printf("json param: %v\n", input)
	b, err := json.Marshal(input)
	if err != nil {
		return fmt.Sprintf("json marchal err %v", err)
	}
	// fmt.Printf("json result: %v\n", b)
	return string(b)
}

func PreloadMake(data interface{}) interface{} {
	rv := reflect.ValueOf(data)
	rt := reflect.TypeOf(data)
	if rt.Kind() == reflect.Map {
		d := data.(map[string]interface{})
		r := make(map[string]interface{})
		for k, v := range d {
			switch v.(type) {
			case string:
				r[k] = preloadParse(v.(string))
				break
			default:
				r[k] = v
				break
			}
		}
		return r
	} else {
		// TODO
		for i := 0; i < rt.NumField(); i += 1 {
			v := rv.Field(i)
			k := v.Kind()
			if k == reflect.String {
				v.SetString("aaa")
			}
		}
	}
	return data
}

func preloadParse(data string) interface{} {
	if strings.HasPrefix(data, "${") && strings.HasSuffix(data, "}") {
		chars := []rune(data[2 : len(data)-1])
		pm := &PreloadParser{
			charIndex:  0,
			charRow:    1,
			charColumn: 1,
			charCount:  len(chars),
			chars:      chars,
		}
		r, err := pm.Parse()
		if err != nil {
			fmt.Printf("parse err: %v\n", err)
			return data
		}
		return r
	} else {
		return data
	}
}
