package xlsrd4

import (
	"fmt"
	"math"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func readUInt2(data []byte, pos int32) (uint16, error) {
	if pos < 0 {
		return 0, fmt.Errorf("无效位置：%d", pos)
	}

	o0 := uint16(data[pos])
	o8 := uint16(data[pos+1]) << 8
	return o0 | o8, nil
}

func readInt4(data []byte, pos int32) (int32, error) {
	size := int32(len(data))
	if pos < 0 || pos > size {
		return 0, fmt.Errorf("读取越界 %d", pos)
	}

	end := int32(math.Min(float64(size), float64(pos)+4))

	var o24 int32
	var o16 int32
	var o8 int32
	var o0 int32

	if end < (pos + 4) {
		o24 = 0
	} else {
		o24 = int32(data[pos+3]) << 24
	}

	if end < (pos + 3) {
		o16 = 0
	} else {
		o16 = int32(data[pos+2]) << 16
	}

	if end < (pos + 2) {
		o8 = 0
	} else {
		o8 = int32(data[pos+1]) << 8
	}

	if end < (pos + 1) {
		o0 = 0
	} else {
		o0 = int32(data[pos])
	}

	return o0 | o8 | o16 | o24, nil
}

// 带1字节长度信息 结果： 字符串数据，带头部的整体长度，错误
func readUnicodeStringShort(subData []byte) ([]byte, uint16, error) {
	charCount := uint16(subData[0])
	v, len, err := ReadUnicodeString(subData[1:], charCount)
	return v, len + 1, err
}

// 带2字节长度信息 结果： 字符串数据，带头部的整体长度，错误
func ReadUnicodeStringLong(subData []byte) ([]byte, uint16, error) {
	charCount, err := readUInt2(subData, 0)
	if err != nil {
		return nil, charCount, err
	}
	v, len, err := ReadUnicodeString(subData[2:], charCount)
	return v, len + 2, err
}

// 结果： 字符串数据，带头部的整体长度，错误
func ReadUnicodeString(subData []byte, chatCount uint16) ([]byte, uint16, error) {
	// bit:0 ; 0 = compression 8bit,  1 = uncompressed 16bit
	isCompressed := (0x01 & subData[0]) == 0

	// bit:2 ;  Asian phonetic settings
	hasAsian := (0x04 & subData[0] >> 2) == 1

	// bit:3 ; Rich-Text settings
	hasRichText := (0x08 & subData[0] >> 3) == 1

	var length uint16
	if isCompressed {
		length = chatCount
	} else {
		length = chatCount * 2
	}

	fmt.Printf("ReadUnicodeString %t %t %t %d\n", isCompressed, hasAsian, hasRichText, length)

	v, err := readUtf8FromUtf16Le(subData[1:1+length], isCompressed)
	return v, length + 1, err
}

func readByteStringStort(code uint16, v []byte) (string, error) {
	ln := v[0]
	return convToUtf8ByCode(code, v[1:1+ln])
}

// 解压，BIFF8 才会使用压缩。
func decompressForBiff8(source []byte) []byte {
	result := make([]byte, 0)
	for i := 0; i < len(source); i += 1 {
		result = append(result, source[i], 0)
	}
	return result
}

func readUtf8FromUtf16Le(source []byte, compressed bool) ([]byte, error) {
	if compressed {
		source = decompressForBiff8(source)
	}
	fmt.Printf("readUtf8FromUtf16Le: %d\n", len(source))
	result, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), source)
	return result, err
}
