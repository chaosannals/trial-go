package xlsrd3

import (
	"bytes"
	"fmt"
	"math"
	"os"
)

var xlsHead []byte

func init() {
	xlsHead = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
}

type xlsFileReader struct {
	xlsBytes []byte
	xlsSize  int32
	pos      int32
}

func (reader *xlsFileReader) traceInfo() {
	fmt.Printf("size: %d\n", len(reader.xlsBytes))
	fmt.Printf("pos: %d\n", reader.pos)
}

func readXlsFile(xlsPath string) (reader *xlsFileReader, err error) {
	reader = &xlsFileReader{
		pos: 0,
	}

	// 读出所有数据
	reader.xlsBytes, err = os.ReadFile(xlsPath)
	if err != nil {
		return
	}
	reader.xlsSize = int32(len(reader.xlsBytes))

	// 验证头
	head := reader.xlsBytes[0:8]
	if !bytes.Equal(xlsHead, head) {
		err = fmt.Errorf("无效的 XLS 头")
		return
	}
	reader.pos = 8

	// 开始读取定位信息
	numBigBlockDepotBlocks, err := reader.readInt4At(NUM_BIG_BLOCK_DEPOT_BLOCKS_POS)
	if err != nil {
		return
	}
	fmt.Printf("全部块数量: %d \n", numBigBlockDepotBlocks)

	rootStartBlock, err := reader.readInt4At(ROOT_START_BLOCK_POS)
	if err != nil {
		return
	}
	fmt.Printf("根开始块ID：%d\n", rootStartBlock)

	sbdStartBlock, err := reader.readInt4At(SMALL_BLOCK_DEPOT_BLOCK_POS)
	if err != nil {
		return
	}
	if sbdStartBlock == -2 {
		fmt.Println("没有 sbd")
	}
	fmt.Printf("sbd ID：%d\n", sbdStartBlock)

	extensionBlock, err := reader.readInt4At(EXTENSION_BLOCK_POS)
	if err != nil {
		return
	}
	if extensionBlock == -2 {
		fmt.Println("没有 extensionBlock")
	}
	fmt.Printf("extensionBlock ID：%d\n", extensionBlock)

	numExtensionBlocks, err := reader.readInt4At(NUM_EXTENSION_BLOCK_POS)
	if err != nil {
		return
	}
	fmt.Printf("numExtensionBlocks ：%d\n", numExtensionBlocks)

	bbdBlocks := numBigBlockDepotBlocks
	if numExtensionBlocks != 0 {
		// 如果启用了扩展区，说明 bb 区被占满了
		bbdBlocks = (BIG_BLOCK_SIZE - BIG_BLOCK_DEPOT_BLOCKS_POS) / 4
	}
	fmt.Printf("bbdBlocks: %d \n", bbdBlocks)

	// TODO READER 如此封装还是不够理想。文件结构导致 pos 并非顺序，而是随机访问
	reader.pos = BIG_BLOCK_DEPOT_BLOCKS_POS
	bigBlockDepotBlocks := make([]int32, bbdBlocks) // 只是索引
	for i := 0; i < int(bbdBlocks); i += 1 {
		bbdPos, err := reader.readInt4()
		if err != nil {
			return nil, err
		}
		fmt.Printf("bbdPos: %d\n", bbdPos)
		bigBlockDepotBlocks[i] = bbdPos
	}

	for i := 0; i < int(numExtensionBlocks); i += 1 {
		reader.pos = int32((extensionBlock + 1) * BIG_BLOCK_SIZE)
		blocksToRead := math.Min(float64(numBigBlockDepotBlocks-bbdBlocks), BIG_BLOCK_SIZE/4-1)

		for i := bbdBlocks; i < bbdBlocks+int32(blocksToRead); i += 1 {
			bbdiPos, err := reader.readInt4()
			if err != nil {
				return nil, err
			}
			fmt.Printf("bbdiPos(%d): %d \n", i, bbdiPos)
		}
	}

	// 开始读数据, TODO 又一个随机乱跳的访问操作
	pos := 0
	bigBlockChain := make([]byte, numBigBlockDepotBlocks*BIG_BLOCK_SIZE)
	// bbs := BIG_BLOCK_SIZE / 4
	for i := 0; i < int(numBigBlockDepotBlocks); i += 1 {
		pos = int(bigBlockDepotBlocks[i]+1) * BIG_BLOCK_SIZE
		copy(bigBlockChain[i*BIG_BLOCK_SIZE:(i+1)*BIG_BLOCK_SIZE], reader.xlsBytes[pos:pos+BIG_BLOCK_SIZE])
		pos += BIG_BLOCK_SIZE
	}
	fmt.Printf("bbc 总大小: %d\n", len(bigBlockChain))

	sbdBlock := sbdStartBlock
	smallBlockChain := make([]byte, 0)
	for sbdBlock != -2 {
		pos = (int(sbdBlock) + 1) * BIG_BLOCK_SIZE
		sbdb := reader.xlsBytes[pos : pos+BIG_BLOCK_SIZE]
		smallBlockChain = append(smallBlockChain, sbdb...)
		pos += BIG_BLOCK_SIZE

		// 这些块 4 应该 int32 的 byte 数，这索引看来，数字有点大，这个定位很魔性。
		r, err := readInt4(bigBlockChain, int32(sbdBlock)*4)
		if err != nil {
			return nil, err
		}
		sbdBlock = r
	}
	fmt.Printf("sbc 总大小: %d\n", len(smallBlockChain))

	entry, err := readEntry(reader.xlsBytes, bigBlockChain, rootStartBlock)
	if err != nil {
		return nil, err
	}
	fmt.Printf("data len: %d\n", len(entry))

	return
}

func readEntry(all []byte, bigBC []byte, block int32) ([]byte, error) {
	result := make([][]byte, 0)
	for block != -2 {
		pos := (block + 1) * BIG_BLOCK_SIZE
		result = append(result, all[pos:pos+BIG_BLOCK_SIZE])
		b, err := readInt4(bigBC, int32(block*4))
		if err != nil {
			return nil, err
		}
		block = b
	}
	bytes := make([]byte, len(result)*BIG_BLOCK_SIZE)
	for i, v := range result {
		copy(bytes[i*BIG_BLOCK_SIZE:(i+1)*BIG_BLOCK_SIZE], v)
	}
	return bytes, nil
}

func (reader *xlsFileReader) readUInt2() (uint16, error) {
	if reader.pos < 0 {
		return 0, fmt.Errorf("无效位置：%d", reader.pos)
	}

	o0 := uint16(reader.xlsBytes[reader.pos])
	o8 := uint16(reader.xlsBytes[reader.pos+1]) << 8
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

func (reader *xlsFileReader) readInt4At(pos int32) (int32, error) {
	return readInt4(reader.xlsBytes, pos)
}

func (reader *xlsFileReader) readInt4() (result int32, err error) {
	result, err = reader.readInt4At(reader.pos)
	reader.pos += 4
	return
}
