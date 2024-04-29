package xlsrd4

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
)

type xlsOleProperty struct {
	Name     string
	TypeId   byte
	Size     int32
	StartPos int32
}

type xlsOlePropertySets struct {
	WorkBookId            int
	RootEntryId           int
	SummaryInfoId         int
	DocumentSummaryInfoId int
	All                   []xlsOleProperty
}

type xlsOleFile struct {
	xlsBytes        []byte
	bigBlockChain   []byte
	smallBlockChain []byte
	propertySet     *xlsOlePropertySets
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

	// 读大区块链
	bigBlockChain, err := readOleBigBlockChain(xlsBytes)
	if err != nil {
		return nil, err
	}
	fmt.Printf("大区块链大小：%d\n", len(bigBlockChain))

	// 读小区块链
	smallBlockChain, err := readOleSmallBlockChain(xlsBytes, bigBlockChain)
	if err != nil {
		return nil, err
	}
	fmt.Printf("小区块链大小：%d\n", len(smallBlockChain))

	// 读根区块
	rootBlock, err := readOleRootBlock(xlsBytes, bigBlockChain)
	if err != nil {
		return nil, err
	}
	fmt.Printf("根区块大小：%d\n", len(rootBlock))

	// 读属性集
	propertySets, err := readOlePropertySets(rootBlock)
	if err != nil {
		return nil, err
	}
	fmt.Printf("属性个数：%d\n", len(propertySets.All))

	return &xlsOleFile{
		xlsBytes:        xlsBytes,
		bigBlockChain:   bigBlockChain,
		smallBlockChain: smallBlockChain,
		propertySet:     propertySets,
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
		d := data[pos : pos+BIG_BLOCK_SIZE]
		smallBlockChain = append(smallBlockChain, d...)
		pos += BIG_BLOCK_SIZE

		r, err := readInt4(bigBlockChain, blockPos*4)
		if err != nil {
			return nil, err
		}
		blockPos = r
	}
	return smallBlockChain, nil
}

func readOleRootBlock(data []byte, bigBlockChain []byte) ([]byte, error) {
	// 头部定位
	startPos, err := readInt4(data, ROOT_START_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	fmt.Printf("根区块ID：%d\n", startPos)

	return readOleBlock(data, bigBlockChain, startPos)
}

func readOlePropertySets(rootBlock []byte) (*xlsOlePropertySets, error) {
	result := &xlsOlePropertySets{
		WorkBookId:            -1,
		RootEntryId:           -1,
		SummaryInfoId:         -1,
		DocumentSummaryInfoId: -1,
		All:                   make([]xlsOleProperty, 0),
	}
	blockSize := len(rootBlock)
	for offset := 0; offset < blockSize; offset += PROPERTY_STORAGE_BLOCK_SIZE {
		d := rootBlock[offset : offset+PROPERTY_STORAGE_BLOCK_SIZE]

		nameSize := int32(d[SIZE_OF_NAME_POS]) | (int32(d[SIZE_OF_NAME_POS+1]) << 8)
		fmt.Printf("属性名大小: %d\n", nameSize)

		typeId := d[TYPE_POS]
		fmt.Printf("属性类型: %d\n", typeId)

		startPos, err := readInt4(d, START_BLOCK_POS)
		if err != nil {
			return nil, err
		}
		fmt.Printf("属性起始位置: %d\n", startPos)

		size, err := readInt4(d, SIZE_POS)
		if err != nil {
			return nil, err
		}
		fmt.Printf("属性大小: %d\n", size)

		name := string(d[0:nameSize])
		fmt.Printf("属性名(%d): %s %v\n", len(name), name, d[0:nameSize])
		tag := convNameTag(d[0:nameSize])
		fmt.Printf("属性名[转换](%d): %s %v\n", len(tag), tag, d[0:nameSize])
		upName := strings.ToUpper(tag)
		fmt.Printf("属性名大写化(%d): %s\n", len(upName), upName)

		result.All = append(result.All, xlsOleProperty{
			Name:     name,
			TypeId:   typeId,
			Size:     size,
			StartPos: startPos,
		})

		// (BIFF5 uses Book, BIFF8 uses Workbook)
		if strings.Compare(upName, "WORKBOOK") == 0 || strings.Compare(upName, "BOOK") == 0 {
			workbook := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("workbook: %d\n", workbook)
			result.WorkBookId = workbook
		} else if upName == "ROOT ENTRY" || upName == "R" {
			rootEntry := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("rootEntry: %d\n", rootEntry)
			result.RootEntryId = rootEntry
		} else if tag == "SummaryInformation" {
			summaryInfo := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("summaryInfo: %d\n", summaryInfo)
			result.SummaryInfoId = summaryInfo
		} else if tag == "DocumentSummaryInformation" {
			docSummaryInfo := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("docSummaryInfo: %d\n", docSummaryInfo)
			result.DocumentSummaryInfoId = docSummaryInfo
		}
	}

	return result, nil
}

func readOleBlock(data []byte, bigBlockChain []byte, blockPos int32) ([]byte, error) {
	result := make([]byte, 0)
	for blockPos != -2 {
		pos := (blockPos + 1) * BIG_BLOCK_SIZE
		d := data[pos : pos+BIG_BLOCK_SIZE]
		result = append(result, d...)

		b, err := readInt4(bigBlockChain, blockPos*4)
		if err != nil {
			return nil, err
		}
		blockPos = b
	}

	return result, nil
}

func readOleStream() {

}

// 这些标签是固定字符集(不确定是否是 Windows 1252，参考 PHPSpreadSheet 做了类似处理)
func convNameTag(data []byte) string {
	result := make([]byte, len(data))
	j := 0
	for i, v := range data {
		if (i == 0 && v == 5) || (v == 0) {
			continue
		}
		result[j] = v
		j += 1
	}
	return string(result[0:j])
}
