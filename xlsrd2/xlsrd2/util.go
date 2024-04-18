package xlsrd2

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
	SMALL_BLOCK_SIZE            = 0x40
	PROPERTY_STORAGE_BLOCK_SIZE = 0x80

	SMALL_BLOCK_THRESHOLD = 0x1000

	REKEY_BLOCK = 0x400

	XLS_TYPE_FORMULA         = 0x0006
	XLS_TYPE_EOF             = 0x000A
	XLS_TYPE_EXTERNSHEET     = 0x0017
	XLS_TYPE_DEFINEDNAME     = 0x0018
	XLS_TYPE_DATEMODE        = 0x0022
	XLS_TYPE_EXTERNNAME      = 0x0023
	XLS_TYPE_FILEPASS        = 0x002F
	XLS_TYPE_FONT            = 0x0031
	XLS_TYPE_CODEPAGE        = 0x0042
	XLS_TYPE_SHEET           = 0x0085
	XLS_TYPE_PALETTE         = 0x0092
	XLS_TYPE_XF              = 0x00E0
	XLS_TYPE_MSODRAWINGGROUP = 0x00EB
	XLS_TYPE_SST             = 0x00FC
	XLS_TYPE_LABELSST        = 0x00FD
	XLS_TYPE_EXTERNALBOOK    = 0x01AE
	XLS_TYPE_NUMBER          = 0x0203
	XLS_TYPE_LABEL           = 0x0204
	XLS_TYPE_BOOLERR         = 0x0205
	XLS_TYPE_RK              = 0x027E
	XLS_TYPE_STYLE           = 0x0293
	XLS_TYPE_FORMAT          = 0x041E
	XLS_TYPE_BOF             = 0x0809
	XLS_TYPE_XFEXT           = 0x087D

	XLS_WORKBOOKGLOBALS = 0x0005
	XLS_WORKSHEET       = 0x0010

	XLS_BIFF8 = 0x0600
	XLS_BIFF7 = 0x0500

	MS_BIFF_CRYPTO_NONE = 0
	MS_BIFF_CRYPTO_XOR  = 1
	MS_BIFF_CRYPTO_RC4  = 2

	WORKSHEET_SHEETSTATE_VISIBLE    = "visible"
	WORKSHEET_SHEETSTATE_HIDDEN     = "hidden"
	WORKSHEET_SHEETSTATE_VERYHIDDEN = "veryHidden"
)

func GetUInt2d(all []byte, pos int) (uint16, error) {
	if pos < 0 {
		return 0, fmt.Errorf("无效位置：%d", pos)
	}

	o0 := uint16(all[pos])
	o8 := uint16(all[pos+1]) << 8
	return o0 | o8, nil
}

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
		// 这么写是因为 PHP 的整形大于 32bit，单纯位移不能使得数为负。golang int32 直接位移即可
		// o24 = int32(all[pos+3])
		// if o24 >= 128 {
		// 	o24 = -int32(math.Abs(float64((256 - o24) << 24)))
		// } else {
		// 	o24 = (o24 & 127) << 24
		// }

		// 与上面的区别不知道是什么。可能是语言问题，golang 应该下面这个就可以。
		o24 = int32(all[pos+3]) << 24
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

func ReadData(all []byte, bigBC []byte, block int32) ([]byte, error) {
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

// 这些标签是固定字符集(不确定是否是 Windows 1252，参考 PHPSpreadSheet 做了类似处理)
func GetNameTag(data []byte) string {
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

func ReadPropertySets(data []byte) (*XlsBookPropertySets, error) {
	result := &XlsBookPropertySets{
		WorkBookId:            nil,
		RootEntryId:           nil,
		SummaryInfoId:         nil,
		DocumentSummaryInfoId: nil,
		All:                   make([]XlsBookProperty, 0),
	}
	dataLen := len(data)
	for offset := 0; offset < dataLen; offset += PROPERTY_STORAGE_BLOCK_SIZE {
		d := data[offset : offset+PROPERTY_STORAGE_BLOCK_SIZE]

		nameSize := int(d[SIZE_OF_NAME_POS]) | (int(d[SIZE_OF_NAME_POS+1]) << 8)

		dType := int(d[TYPE_POS])
		fmt.Printf("d type: %d\n", dType)

		startBlock, err := GetInt4d(d, START_BLOCK_POS)
		if err != nil {
			return nil, err
		}
		fmt.Printf("d startBlock: %d\n", startBlock)

		size, err := GetInt4d(d, SIZE_POS)
		if err != nil {
			return nil, err
		}
		fmt.Printf("d size: %d\n", size)

		name := string(d[0:nameSize])
		fmt.Printf("name(%d): %s %v\n", len(name), name, d[0:nameSize])
		tag := GetNameTag(d[0:nameSize])
		upName := strings.ToUpper(tag)
		fmt.Printf("up name(%d): %s\n", len(upName), upName)

		result.All = append(result.All, XlsBookProperty{
			Name:     name,
			Size:     size,
			TypeId:   dType,
			StartPos: startBlock,
		})

		// (BIFF5 uses Book, BIFF8 uses Workbook)
		if strings.Compare(upName, "WORKBOOK") == 0 || strings.Compare(upName, "BOOK") == 0 {
			workbook := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("workbook: %d\n", workbook)
			result.WorkBookId = &workbook
		} else if upName == "ROOT ENTRY" || upName == "R" {
			rootentry := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("rootentry: %d\n", rootentry)
			result.RootEntryId = &rootentry
		} else if tag == "SummaryInformation" {
			summaryInfo := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("summaryInfo: %d\n", summaryInfo)
			result.SummaryInfoId = &summaryInfo
		} else if tag == "DocumentSummaryInformation" {
			docSummaryInfo := offset / PROPERTY_STORAGE_BLOCK_SIZE
			fmt.Printf("docSummaryInfo: %d\n", docSummaryInfo)
			result.DocumentSummaryInfoId = &docSummaryInfo
		}
	}
	return result, nil
}

func ReadXls(xlsPath string) (*XlsBook, error) {
	xlsBytes, err := os.ReadFile(xlsPath)
	if err != nil {
		return nil, err
	}
	xlsHead := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
	head := xlsBytes[0:8]
	if bytes.Equal(xlsHead, head) {
		fmt.Println("有效 xls 文件头")
	} else {
		return nil, fmt.Errorf("无效 xls 文件头 %s \n", hex.Dump(head))
	}

	// 开始读取定位信息
	//
	numBigBlockDepotBlocks, err := GetInt4d(xlsBytes, NUM_BIG_BLOCK_DEPOT_BLOCKS_POS)
	if err != nil {
		return nil, err
	}
	fmt.Printf("全部块数量: %d \n", numBigBlockDepotBlocks)

	rootStartBlock, err := GetInt4d(xlsBytes, ROOT_START_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	fmt.Printf("根开始块ID：%d\n", rootStartBlock)

	sbdStartBlock, err := GetInt4d(xlsBytes, SMALL_BLOCK_DEPOT_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	if sbdStartBlock == -2 {
		fmt.Println("没有 sbd")
	}
	fmt.Printf("sbd ID：%d\n", sbdStartBlock)

	extensionBlock, err := GetInt4d(xlsBytes, EXTENSION_BLOCK_POS)
	if err != nil {
		return nil, err
	}
	if extensionBlock == -2 {
		fmt.Println("没有 extensionBlock")
	}
	fmt.Printf("extensionBlock ID：%d\n", extensionBlock)

	numExtensionBlocks, err := GetInt4d(xlsBytes, NUM_EXTENSION_BLOCK_POS)
	if err != nil {
		return nil, err
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
			return nil, err
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
				return nil, err
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
	smallBlockChain := make([]byte, 0)
	for sbdBlock != -2 {
		pos = (int(sbdBlock) + 1) * BIG_BLOCK_SIZE
		sbdb := xlsBytes[pos : pos+BIG_BLOCK_SIZE]
		smallBlockChain = append(smallBlockChain, sbdb...)
		pos += BIG_BLOCK_SIZE

		// TODO 这里 sbdBlock * 4 不理解。
		r, err := GetInt4d(bigBlockChain, int(sbdBlock)*4)
		if err != nil {
			return nil, err
		}
		sbdBlock = r
	}

	entry, err := ReadData(xlsBytes, bigBlockChain, rootStartBlock)
	if err != nil {
		return nil, err
	}
	fmt.Printf("data len: %d\n", len(entry))

	ps, err := ReadPropertySets(entry)
	if err != nil {
		return nil, err
	}

	book := &XlsBook{
		bigBlockChain:   bigBlockChain,
		smallBlockChain: smallBlockChain,
		PropertySets:    ps,
		entry:           entry,
		xlsBytes:        xlsBytes,
	}
	bd, err := book.ReadStream(int32(*book.PropertySets.WorkBookId))
	if err != nil {
		return nil, err
	}
	book.Data = bd

	bsi, err := book.ReadStream((int32(*book.PropertySets.SummaryInfoId)))
	if err != nil {
		return nil, err
	}
	book.SummaryInfo = bsi

	bdsi, err := book.ReadStream(int32(*book.PropertySets.DocumentSummaryInfoId))
	if err != nil {
		return nil, err
	}
	book.DocumentSummaryInfo = bdsi

	return book, nil
}

func GetUtf8FromUtf16Le(source []byte, compressed bool) ([]byte, error) {
	if compressed {
		source = Uncompress(source)
	}
	fmt.Printf("GetUtf8FromUtf16Le: %d\n", len(source))
	result, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), source)
	return result, err
}

// 解压，BIFF8 才会使用压缩。
func Uncompress(source []byte) []byte {
	result := make([]byte, 0)
	for i := 0; i < len(source); i += 1 {
		result = append(result, source[i], 0)
	}
	return result
}

// 带1字节长度信息 结果： 字符串数据，带头部的整体长度，错误
func ReadUnicodeStringShort(subData []byte) ([]byte, int, error) {
	charCount := int(subData[0])
	v, len, err := ReadUnicodeString(subData[1:], charCount)
	return v, len + 1, err
}

// 带2字节长度信息 结果： 字符串数据，带头部的整体长度，错误
func ReadUnicodeStringLong(subData []byte) ([]byte, int, error) {
	charCount, err := GetUInt2d(subData, 0)
	if err != nil {
		return nil, int(charCount), err
	}
	v, len, err := ReadUnicodeString(subData[2:], int(charCount))
	return v, len + 2, err
}

// 结果： 字符串数据，带头部的整体长度，错误
func ReadUnicodeString(subData []byte, chatCount int) ([]byte, int, error) {
	// bit:0 ; 0 = compression 8bit,  1 = uncompressed 16bit
	isCompressed := (0x01 & subData[0]) == 0

	// bit:2 ;  Asian phonetic settings
	hasAsian := (0x04 & subData[0] >> 2) == 1

	// bit:3 ; Rich-Text settings
	hasRichText := (0x08 & subData[0] >> 3) == 1

	var length int
	if isCompressed {
		length = chatCount
	} else {
		length = chatCount * 2
	}

	fmt.Printf("ReadUnicodeString %t %t %t %d\n", isCompressed, hasAsian, hasRichText, length)

	v, err := GetUtf8FromUtf16Le(subData[1:1+length], isCompressed)
	return v, length + 1, err
}

func ReadByteStringStort(code uint16, v []byte) (string, error) {
	ln := v[0]
	return ConvToUtf8ByCode(code, v[1:1+ln])
}

var ColumnIndexStringCache map[uint16]string

const ColumnIndexStringLookup = " ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ColumnIndexToString(index uint16) string {
	if ColumnIndexStringCache == nil {
		ColumnIndexStringCache = map[uint16]string{}
	}

	r, isOk := ColumnIndexStringCache[index]
	if isOk {
		return r
	}

	indexValue := index
	result := ""
	for {
		c := (indexValue % 26)
		if c == 0 {
			c = 26
		}
		indexValue = (indexValue - c) / 26
		v := ColumnIndexStringLookup[c : c+1]
		result = fmt.Sprintf("%s%s", v, result)
		if indexValue <= 0 {
			break
		}
	}
	ColumnIndexStringCache[index] = result

	return result
}
