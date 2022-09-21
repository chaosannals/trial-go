package stress

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var randCharset = []rune("1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
var randMaxchar = len(randCharset) - 1

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
	c := self.peekChar()
	if unicode.IsDigit(c) {
		return self.matchNumber()
	}
	if c == '"' {
		return self.matchString()
	}
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

func (self *PreloadParser) matchNumber() (float64, error) {
	buf := make([]rune, 0, 128)
	for {
		c := self.peekChar()
		if c == '.' || unicode.IsDigit(c) {
			self.popChar()
			buf = append(buf, c)
		} else {
			break
		}
	}
	return strconv.ParseFloat(string(buf), 64)
}

func (self *PreloadParser) matchString() (string, error) {
	buf := make([]rune, 0, 128)
	for {
		c := self.peekChar()
		if '"' == c {
			break
		}
		self.popChar()
		buf = append(buf, c)
	}
	return string(buf[1 : len(buf)-1]), nil
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
	//fmt.Println("randomObject start")
	fc := 6 + rand.Intn(20)
	for i := 0; i < fc; i++ {
		fn := self.RandomString(nil)
		r := rand.Intn(100)
		switch {
		case r == 1:
			result[fn] = self.RandomObject(nil)
			break
		case r > 1 && r <= 10:
			result[fn] = rand.Float32()
			break
		case r > 10 && r <= 20:
			result[fn] = self.RandomString(nil)
			break
		default:
			result[fn] = rand.Intn(1000)
			break
		}
	}
	//fmt.Println("randomObject end")
	return result
}

func (self *PreloadParser) RandomList(_ interface{}) interface{} {
	lc := 20 + rand.Intn(100)
	buf := make([]interface{}, 0, lc)
	for i := 0; i < lc; i++ {
		buf = append(buf, self.RandomObject(nil))
	}
	return buf
}

func (self *PreloadParser) RandomString(_ interface{}) string {
	sb := strings.Builder{}
	sn := rand.Intn(5) + 3
	sb.Grow(sn)
	for i := 0; i < sn; i++ {
		ci := rand.Intn(randMaxchar)
		sb.WriteRune(randCharset[ci])
	}
	return sb.String()
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
	switch rt.Kind() {
	case reflect.Map:
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
	case reflect.String:
		return preloadParse(data.(string))
	default:
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
