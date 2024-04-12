package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
)

const (
	// 这些相对于整个文件偏移量
	BIG_BLOCK_DEPOT_BLOCKS_POS     = 0x4C
	NUM_BIG_BLOCK_DEPOT_BLOCKS_POS = 0x2C // 块总量位置
	ROOT_START_BLOCK_POS           = 0x30 // 根开始块ID 位置
	SMALL_BLOCK_DEPOT_BLOCK_POS    = 0x3C //
	EXTENSION_BLOCK_POS            = 0x44
	NUM_EXTENSION_BLOCK_POS        = 0x48

	// 这些是相对与 data 各个块内 的偏移量
	SIZE_OF_NAME_POS = 0x40
	TYPE_POS         = 0x42
	START_BLOCK_POS  = 0x74
	SIZE_POS         = 0x78

	BIG_BLOCK_SIZE              = 0x200
	PROPERTY_STORAGE_BLOCK_SIZE = 0x80
)

func GetInt4d(all []byte, pos int) (int32, error) {
	if pos < 0 {
		return 0, fmt.Errorf("无效位置：%d", pos)
	}

	len := len(all)

	var o24 int32
	var o16 int32
	var o8 int32
	var o0 int32
	if len < (pos + 4) {
		o24 = 0
	} else {
		o24 = int32(all[pos+3])
		if o24 >= 128 {
			o24 = -int32(math.Abs(float64((256 - o24) << 24)))
		} else {
			o24 = (o24 & 127) << 24
		}

		// 与上面的区别不知道是什么。可能是语言问题，golang 应该下面这个就可以。
		// o24 = int32(all[pos+3]) << 24
	}

	if len < (pos + 3) {
		o16 = 0
	} else {
		o16 = int32(all[pos+2]) << 16
	}

	if len < (pos + 2) {
		o8 = 0
	} else {
		o8 = int32(all[pos+1]) << 8
	}

	if len < (pos + 1) {
		o0 = 0
	} else {
		o0 = int32(all[pos])
	}

	return o0 | o8 | o16 | o24, nil
}

func readData(all []byte, bigBC []byte, block int32) ([]byte, error) {
	result := make([][]byte, 0)
	for block != -2 {
		pos := (block + 1) * BIG_BLOCK_SIZE
		result = append(result, all[pos:pos+BIG_BLOCK_SIZE])

		// TODO 这里取下一个也是 block * 4 不理解
		b, err := GetInt4d(bigBC, int(block*4))
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

func readPropertySets(data []byte) error {

	dataLen := len(data)
	for offset := 0; offset < dataLen; {
		d := data[offset : offset+PROPERTY_STORAGE_BLOCK_SIZE]

		nameSize := int(d[SIZE_OF_NAME_POS]) | (int(d[SIZE_OF_NAME_POS+1]) << 8)

		dType := int(d[TYPE_POS])
		fmt.Printf("d type: %d\n", dType)

		startBlock, err := GetInt4d(d, START_BLOCK_POS)
		if err != nil {
			return err
		}
		fmt.Printf("d startBlock: %d\n", startBlock)

		size, err := GetInt4d(d, SIZE_POS)
		if err != nil {
			return err
		}
		fmt.Printf("d size: %d\n", size)

		name := d[0:nameSize]
		fmt.Printf("d name: %s\n", name)

		offset += PROPERTY_STORAGE_BLOCK_SIZE
	}
	return nil
}

func main() {
	wkDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wkDir)
	xlsPath := filepath.Join(wkDir, "b.xls")
	fmt.Println(xlsPath)
	xlsBytes, err := os.ReadFile(xlsPath)
	if err != nil {
		log.Fatal(err)
	}
	head := xlsBytes[0:8]

	// 校验头
	xlsHead := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}

	if bytes.Equal(xlsHead, head) {
		fmt.Println("有效 xls 文件头")
	} else {
		fmt.Printf("无效 xls 文件头 %s \n", hex.Dump(head))
	}

	// 开始读取定位信息
	//
	numBigBlockDepotBlocks, err := GetInt4d(xlsBytes, NUM_BIG_BLOCK_DEPOT_BLOCKS_POS)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("全部块数量: %d \n", numBigBlockDepotBlocks)

	rootStartBlock, err := GetInt4d(xlsBytes, ROOT_START_BLOCK_POS)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("根开始块ID：%d\n", rootStartBlock)

	sbdStartBlock, err := GetInt4d(xlsBytes, SMALL_BLOCK_DEPOT_BLOCK_POS)
	if err != nil {
		log.Fatal(err)
	}
	if sbdStartBlock == -2 {
		fmt.Println("没有 sbd")
	}
	fmt.Printf("sbd ID：%d\n", sbdStartBlock)

	extensionBlock, err := GetInt4d(xlsBytes, EXTENSION_BLOCK_POS)
	if err != nil {
		log.Fatal(err)
	}
	if extensionBlock == -2 {
		fmt.Println("没有 extensionBlock")
	}
	fmt.Printf("extensionBlock ID：%d\n", extensionBlock)

	numExtensionBlocks, err := GetInt4d(xlsBytes, NUM_EXTENSION_BLOCK_POS)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("numExtensionBlocks ：%d\n", numExtensionBlocks)

	bbdBlocks := numBigBlockDepotBlocks
	if numExtensionBlocks != 0 {
		// 如果启用了扩展区，说明 bb 区被占满了
		bbdBlocks = (BIG_BLOCK_SIZE - BIG_BLOCK_DEPOT_BLOCKS_POS) / 4
	}
	fmt.Printf("bbdBlocks: %d \n", bbdBlocks)

	pos := BIG_BLOCK_DEPOT_BLOCKS_POS
	bigBlockDepotBlocks := make([]int32, bbdBlocks) // 只是索引
	for i := 0; i < int(bbdBlocks); i += 1 {
		bbdPos, err := GetInt4d(xlsBytes, pos)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("bbdPos: %d\n", bbdPos)
		bigBlockDepotBlocks[i] = bbdPos
		pos += 4
	}

	for i := 0; i < int(numExtensionBlocks); i += 1 {
		pos = int((extensionBlock + 1) * BIG_BLOCK_SIZE)
		blocksToRead := math.Min(float64(numBigBlockDepotBlocks-bbdBlocks), BIG_BLOCK_SIZE/4-1)

		for i := bbdBlocks; i < bbdBlocks+int32(blocksToRead); i += 1 {
			bbdiPos, err := GetInt4d(xlsBytes, pos)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("bbdiPos(%d): %d \n", i, bbdiPos)
			pos += 4
		}
	}

	// 开始读数据
	pos = 0
	bigBlockChain := make([]byte, numBigBlockDepotBlocks*BIG_BLOCK_SIZE)
	// bbs := BIG_BLOCK_SIZE / 4
	for i := 0; i < int(numBigBlockDepotBlocks); i += 1 {
		pos = int(bigBlockDepotBlocks[i]+1) * BIG_BLOCK_SIZE
		copy(bigBlockChain[i*BIG_BLOCK_SIZE:(i+1)*BIG_BLOCK_SIZE], xlsBytes[pos:pos+BIG_BLOCK_SIZE])
		pos += BIG_BLOCK_SIZE
	}
	fmt.Printf("bbc 总大小: %d\n", len(bigBlockChain))

	sbdBlock := sbdStartBlock
	smallBlockChain := make([][]byte, 0)
	for sbdBlock != -2 {
		pos = (int(sbdBlock) + 1) * BIG_BLOCK_SIZE
		sbdb := xlsBytes[pos : pos+BIG_BLOCK_SIZE]
		smallBlockChain = append(smallBlockChain, sbdb)
		pos += BIG_BLOCK_SIZE

		// TODO 这里 sbdBlock * 4 不理解。
		r, err := GetInt4d(bigBlockChain, int(sbdBlock)*4)
		if err != nil {
			log.Fatal(err)
		}
		sbdBlock = r
	}

	entry, err := readData(xlsBytes, bigBlockChain, rootStartBlock)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("data len: %d\n", len(entry))

	err = readPropertySets(entry)
	if err != nil {
		log.Fatal(err)
	}
}
