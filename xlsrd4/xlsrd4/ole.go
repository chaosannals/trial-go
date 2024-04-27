package xlsrd4

import (
	"bytes"
	"fmt"
	"os"
)

type xlsOleFile struct {
	xlsBytes []byte
}

var xlsHead []byte

func init() {
	xlsHead = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
}

func readOleFile(xlsPath string) (*xlsOleFile, error) {
	xlsBytes, err := os.ReadFile(xlsPath)
	if err != nil {
		return nil, err
	}

	// 验证头
	head := xlsBytes[0:8]
	if !bytes.Equal(xlsHead, head) {
		return nil, fmt.Errorf("无效的 XLS 头")
	}

	// 开始读取定位信息
	numBigBlockDepotBlocks, err := readInt4(xlsBytes, NUM_BIG_BLOCK_DEPOT_BLOCKS_POS)
	if err != nil {
		return nil, err
	}
	fmt.Printf("全部块数量: %d \n", numBigBlockDepotBlocks)

	return &xlsOleFile{
		xlsBytes: xlsBytes,
	}, nil
}

func readOleBlock() {

}

func readOleStream() {

}
