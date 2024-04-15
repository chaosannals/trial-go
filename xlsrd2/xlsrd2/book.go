package xlsrd2

import "fmt"

type XlsBookProperty struct {
	TypeId   int
	StartPos int32
	Name     string
	Size     int32
}

type XlsBookPropertySets struct {
	WorkBookId            *int
	RootEntryId           *int
	SummaryInfoId         *int
	DocumentSummaryInfoId *int
	All                   []XlsBookProperty
}

type XlsBook struct {
	// bigBlockDepotBlockCount     int
	// rootStartBlockId            int
	// smallBlockDepotStartBlockId int
	// extensionBlockId            int
	// extensionBlockCount         int

	bigBlockChain       []byte
	smallBlockChain     []byte
	entry               []byte
	xlsBytes            []byte
	PropertySets        *XlsBookPropertySets
	Data                []byte
	SummaryInfo         []byte
	DocumentSummaryInfo []byte
}

func (i *XlsBook) ReadStream(id int32) ([]byte, error) {
	streamData := make([]byte, 0)
	if id < 0 {
		return streamData, fmt.Errorf("不是有效的 Sheet ID")
	}

	size := i.PropertySets.All[id].Size
	startPos := i.PropertySets.All[id].StartPos

	// 如果没有超过小块阈值，就按小块处理
	if size < SMALL_BLOCK_THRESHOLD {
		fmt.Printf("stream is small block: %d\n", id)
		rootData, err := ReadData(i.xlsBytes, i.bigBlockChain, i.PropertySets.All[*i.PropertySets.RootEntryId].StartPos)
		if err != nil {
			return streamData, err
		}
		block := startPos
		for block != -2 {
			// 小块是从 0 块开始的。
			pos := block * SMALL_BLOCK_SIZE
			// 小块是从 rootEntry 块获取数据
			streamData = append(streamData, rootData[pos:pos+SMALL_BLOCK_SIZE]...)
			b, err := GetInt4d(i.smallBlockChain, int(block)*4)
			if err != nil {
				return streamData, err
			}
			block = b
		}

		return streamData, nil
	}

	// 不是小块就是大块, 向上取整块数
	fmt.Printf("stream is big block: %d\n", id)
	numBlocks := size / BIG_BLOCK_SIZE
	if (size % BIG_BLOCK_SIZE) != 0 {
		numBlocks += 1
	}

	// 空快
	if numBlocks == 0 {
		return streamData, nil
	}

	block := startPos
	for block != -2 {
		// 大块是从 1 块开始的
		pos := (block + 1) * BIG_BLOCK_SIZE
		// 大块是从 i.xlsBytes 整体上获取数据
		streamData = append(streamData, i.xlsBytes[pos:pos+BIG_BLOCK_SIZE]...)
		b, err := GetInt4d(i.bigBlockChain, int(block)*4)
		if err != nil {
			return streamData, err
		}
		block = b
	}
	return streamData, nil
}

// func (i *XlsBook) ListWookSheetNames() ([]string, error) {
// 	dLen := len(i.data)
// 	pos := 0
// 	sheets := make([]string, 0)

// 	for ;pos < dLen; {
// 		code, err := GetUInt2d(i.data, pos)
// 		if err != nil {
// 			return sheets, err
// 		}

// 	}
// }
