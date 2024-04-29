package xlsrd4

import (
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var iconvEncodings map[uint16]encoding.Encoding

func init() {
	iconvEncodings = map[uint16]encoding.Encoding{
		0:  charmap.Windows1252,
		37: charmap.CodePage037,
		// 367: charmap.CodePage037, // ASCII ,不太确定是不是这个 037
		367: unicode.UTF8, // UTF8 完全兼容 ASCII
		437: charmap.CodePage437,
		850: charmap.CodePage850,
		852: charmap.CodePage852,
		855: charmap.CodePage855,
		858: charmap.CodePage858,
		860: charmap.CodePage860,
		862: charmap.CodePage862,
		863: charmap.CodePage863,
		865: charmap.CodePage865,
		866: charmap.CodePage866,
		874: charmap.Windows874, // ANSI Thai
		932: japanese.ShiftJIS,  //ANSI Japanese Shift-JIS
		936: simplifiedchinese.GBK,
		949: korean.EUCKR, // ANSI Korean (Wansung)
		950: traditionalchinese.Big5,

		1200: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), //    UTF-16 (BIFF8)
		1250: charmap.Windows1250,

		10000: charmap.Macintosh,          //    Apple Roman
		10001: japanese.ShiftJIS,          //    Macintosh Japanese
		10002: traditionalchinese.Big5,    //    Macintosh Chinese Traditional
		10007: charmap.MacintoshCyrillic,  //    Macintosh Cyrillic
		10008: simplifiedchinese.HZGB2312, //    Macintosh - Simplified Chinese (GB 2312)

		21010: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), //    UTF-16 (BIFF8) This isn"t correct, but some Excel writer libraries erroneously use Codepage 21010 for UTF-16LE

		// 65000: unicode.UTF7,
		65001: unicode.UTF8,
	}
}

func convToUtf8ByCode(code uint16, data []byte) (string, error) {
	coding, isOk := iconvEncodings[code]
	if isOk {
		r, _, err := transform.Bytes(coding.NewDecoder(), data)
		if err != nil {
			return "", err
		}
		return string(r), nil
	}
	return "", fmt.Errorf("unsupported code: %d", code)
}
