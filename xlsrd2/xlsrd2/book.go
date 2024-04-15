package xlsrd2

import (
	"fmt"
)

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
	xlsBytes            []byte // 对应 OLERead 类 data
	PropertySets        *XlsBookPropertySets
	Data                []byte // 对应 xls 类 data
	SummaryInfo         []byte
	DocumentSummaryInfo []byte
	version int16
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

func (i *XlsBook) ListWookSheetNames() ([]string, error) {
	dLen := len(i.Data)
	pos := 0
	sheets := make([]string, 0)

	for pos < dLen {
		code, err := GetUInt2d(i.Data, pos)
		if err != nil {
			return sheets, err
		}
		switch code {
		case XLS_TYPE_BOF:
			p, err := i.ReadBof(pos)
			if err != nil {
				return sheets, err
			}
			pos = p
			break
		case XLS_TYPE_SHEET:
			p, err := i.ReadSheet(pos)
			if err != nil {
				return sheets, err
			}
			pos = p
			break
		}
	}
}

func (i *XlsBook) ReadBof(pos int) (int, error) {
	length, err := GetUInt2d(i.Data, pos+2)
	if err != nil {
		return pos, err
	}

	recordData := i.Data[pos+4 : length]
	pos += 4 + int(length)
	substreamType, err := GetUInt2d(recordData, 2)
	if err != nil {
		return pos, err
	}

	switch substreamType {
	case XLS_WORKBOOKGLOBALS:
		version, err := GetUInt2d(recordData, 0)
		if err != nil {
			return pos, err
		}
		if (version != XLS_BIFF8) && (version != XLS_BIFF7) {
			return pos, fmt.Errorf("这个文件的格式太老了")
		}
		// TODO 这个结构要读到这一步才能获取到版本信息，加载流程确认下
		fmt.Printf("xls 文件版本: %d", version)
		i.version = version
		break
	case XLS_WORKSHEET:
		// 此项的版本信息不可靠
		// (OpenOffice doc, 5.8)指出使用全局流里面的版本信息
		break
	default:
		for {
			code, err := GetUInt2d(i.Data, pos)
			if err != nil {
				return pos, err
			}
			p, err := i.ReadDefault()
			if err != nil {
				return p, err
			}
			pos = p
			if code == XLS_TYPE_EOF || pos >= len(i.Data) {
				break
			}
		}
		break
	}
	return pos, nil
}

func (i *XlsBook) ReadSheet(pos int)  (int, error) {
	length, err := GetUInt2d(i.Data, pos+2)
	if err != nil {
		return pos, err
	}
	recordData := i.Data[pos+4 : length]

	recOffset, err := GetInt4d(i.Data, pos + 4)
	if err != nil {
		return pos, err
	}
	var sheetState string
	switch recordData[4] {
	case 0x00:
		sheetState = WORKSHEET_SHEETSTATE_VISIBLE
		break
	case 0x01:
		sheetState = WORKSHEET_SHEETSTATE_HIDDEN
		break
	case 0x02:
		sheetState = WORKSHEET_SHEETSTATE_VERYHIDDEN
		break
	default:
		sheetState = WORKSHEET_SHEETSTATE_VISIBLE
		break
	}
	fmt.Printf("sheet state: %s", sheetState)

	sheetType := recordData[5]
	fmt.Printf("sheet type: %d", sheetType)

	recName := ""
	switch i.version {
	case XLS_BIFF8:
		recName = ""
	}

}

func (i *XlsBook) ReadDefault() (int, error) {

}
