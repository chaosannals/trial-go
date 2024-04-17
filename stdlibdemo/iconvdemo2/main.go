package main

import (
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var iconvEncodings map[uint16]encoding.Encoding
var testTexts []string

func init() {
	testTexts = []string{
		"1234567890ABCDEF",
		"1234567890中文ABCDEF",
	}
	iconvEncodings = map[uint16]encoding.Encoding{
		// 367: ascii.
		437:   charmap.CodePage437,
		1200:  unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), //    UTF-16 (BIFF8)
		21010: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), //    UTF-16 (BIFF8) This isn"t correct, but some Excel writer libraries erroneously use Codepage 21010 for UTF-16LE
	}

}

func main() {
	for i, text := range testTexts {
		fmt.Printf("(%d)=====================================\n", i)
		fmt.Printf("source text: %s\n", text)
		for codePage := range iconvEncodings {
			encoding := iconvEncodings[codePage]
			encoder := encoding.NewEncoder()
			eBytes, n, err := transform.Bytes(encoder, []byte(text))
			if err != nil {
				fmt.Printf("encode (%d) error: %v", codePage, err)
				continue
			}
			eText := string(eBytes)
			fmt.Printf("encoding codePage(%d): (%d, %d) => %s\n", codePage, len(eBytes), n, eText)
			decoder := encoding.NewDecoder()
			dBytes, n, err := transform.Bytes(decoder, eBytes)
			if err != nil {
				fmt.Printf("decode (%d) error: %v", codePage, err)
				continue
			}
			dText := string(dBytes)
			fmt.Printf("decoding codePage(%d): (%d, %d) => %s\n", codePage, len(dBytes), n, dText)
		}
	}
}
