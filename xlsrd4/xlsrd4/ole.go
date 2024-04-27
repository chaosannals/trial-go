package xlsrd4

import (
	"bytes"
	"fmt"
	"math"
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

	// 头部定位
	rootStartBlock, err := readInt4(xlsBytes, ROOT_START_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	fmt.Printf("根开始块ID：%d\n", rootStartBlock)

	// 大区块链
	bigBlockChain, err := readOleBigBlockChain(xlsBytes)
	if err != nil {
		return nil, err
	}
	fmt.Printf("大区块链大小：%d\n", len(bigBlockChain))

	// 小区块链
	smallBlockChain, err := readOleSmallBlockChain(xlsBytes, bigBlockChain)
	if err != nil {
		return nil, err
	}
	fmt.Printf("小区块链大小：%d\n", len(smallBlockChain))

	return &xlsOleFile{
		xlsBytes: xlsBytes,
	}, nil
}

func readOleBigBlockChain(data []byte) ([]byte, error) {
	// 大区块的数量，总体的数量如果有扩展区，是会算上扩展区的。
	depotCount, err := readInt4(data, NUM_BIG_BLOCK_DEPOT_BLOCKS_POS)
	if err != nil {
		return nil, err
	}
	fmt.Printf("全部大区块数量: %d \n", depotCount)

	// 扩展区定位
	extensionPos, err := readInt4(data, EXTENSION_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	if extensionPos == -2 {
		fmt.Println("没有扩展区")
	}
	fmt.Printf("扩展区定位 ID：%d\n", extensionPos)
	// 扩展区个数
	extensionCount, err := readInt4(data, NUM_EXTENSION_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	fmt.Printf("扩展区数量 ：%d\n", extensionCount)

	// 大区块定位
	pos := int32(BIG_BLOCK_DEPOT_BLOCKS_POS)
	bigBlockDepotBlocks := make([]int32, depotCount) // 只是索引

	// 具体在大区块区内的数量，不会超过上线，超过部分的块，放到扩展区。
	bigBlockDepotCount := depotCount
	if extensionCount != 0 {
		// 如果启用了扩展区，说明 大区块区被占满了
		bigBlockDepotCount = (BIG_BLOCK_SIZE - BIG_BLOCK_DEPOT_BLOCKS_POS) / 4
	}
	fmt.Printf("大区块区内块数量: %d \n", bigBlockDepotCount)

	// 读取大区块内部的定位信息
	for i := 0; i < int(bigBlockDepotCount); i += 1 {
		bbdPos, err := readInt4(data, pos)
		if err != nil {
			return nil, err
		}
		fmt.Printf("大区块: %d\n", bbdPos)
		bigBlockDepotBlocks[i] = bbdPos
		pos += 4
	}

	// 读取扩展区内部的块定位信息
	for i := 0; i < int(extensionCount); i += 1 {
		pos = int32((extensionPos + 1) * BIG_BLOCK_SIZE)
		blockCount := int32(math.Min(float64(depotCount-bigBlockDepotCount), BIG_BLOCK_SIZE/4-1))
		blockEnd := bigBlockDepotCount + blockCount

		for i := bigBlockDepotCount; i < blockEnd; i += 1 {
			bbdiPos, err := readInt4(data, pos)
			if err != nil {
				return nil, err
			}
			bigBlockDepotBlocks[i] = bbdiPos
			fmt.Printf("扩展块(%d): %d \n", i, bbdiPos)
			pos += 4
		}

		bigBlockDepotCount += blockCount
		if bigBlockDepotCount < depotCount {
			nextPos, err := readInt4(data, pos)
			if err != nil {
				return nil, err
			}
			extensionPos = nextPos
		}
	}

	// 读取大区块内容
	pos = 0
	bigBlockChain := make([]byte, depotCount*BIG_BLOCK_SIZE)
	// bbs := BIG_BLOCK_SIZE / 4
	for i := 0; i < int(depotCount); i += 1 {
		pos = int32(bigBlockDepotBlocks[i]+1) * BIG_BLOCK_SIZE
		copy(bigBlockChain[i*BIG_BLOCK_SIZE:(i+1)*BIG_BLOCK_SIZE], data[pos:pos+BIG_BLOCK_SIZE])
		pos += BIG_BLOCK_SIZE
	}
	return bigBlockChain, nil
}

func readOleSmallBlockChain(data []byte, bigBlockChain []byte) ([]byte, error) {
	startPos, err := readInt4(data, SMALL_BLOCK_DEPOT_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	if startPos == -2 {
		fmt.Println("没有小区块")
	}
	fmt.Printf("小区块 ID：%d\n", startPos)

	// 读取数据
	blockPos := startPos
	smallBlockChain := make([]byte, 0)
	for blockPos != -2 {
		pos := (blockPos + 1) * BIG_BLOCK_SIZE
		sbdb := data[pos : pos+BIG_BLOCK_SIZE]
		smallBlockChain = append(smallBlockChain, sbdb...)
		pos += BIG_BLOCK_SIZE

		r, err := readInt4(bigBlockChain, blockPos*4)
		if err != nil {
			return nil, err
		}
		blockPos = r
	}
	return smallBlockChain, nil
}

func readOleBlock() {

}

func readOleStream() {

}
