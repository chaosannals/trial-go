package xlsrd4

import (
	"fmt"
	"math"
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
