package xlsrd2

import (
	"crypto/rc4"
	"fmt"
	"math"
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

type XlsBookSheet struct {
	Name   string
	Offset int32
	State  string
	Type   uint8
}

type XlsBookWorkSheet struct {
	Name             string
	LastColumnLetter string
	LastColumnIdnex  uint16
	TotalRows        uint16
	TotalColumns     uint16
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
	version             uint16

	// TODO 加密部分不可靠
	encryption         int
	encryptionStartPos int32
	rc4Pos             int
	rc4Cipher          *rc4.Cipher

	codePageCode uint16
	codePage     string

	Sheets []XlsBookSheet
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

func (i *XlsBook) ListWookSheetInfos() ([]*XlsBookWorkSheet, error) {
	dLen := len(i.Data)
	pos := 0
	sheets := make([]*XlsBookWorkSheet, 0)

	fmt.Printf("ListWookSheetNames: %d\n", dLen)

Loop1:
	for pos < dLen {
		// fmt.Printf("ListWookSheetNames loop: %d\n", pos)
		code, err := GetUInt2d(i.Data, pos)
		if err != nil {
			return sheets, err
		}
		// fmt.Printf("ListWookSheetNames code: %d %d\n", pos, code)
		switch code {
		case XLS_TYPE_BOF:
			p, err := i.ReadBof(pos)
			if err != nil {
				return sheets, err
			}
			pos = p
		case XLS_TYPE_SHEET:
			p, err := i.ReadSheet(pos)
			if err != nil {
				return sheets, err
			}
			pos = p
		case XLS_TYPE_EOF:
			p, err := i.ReadDefault(pos)
			if err != nil {
				return sheets, err
			}
			pos = p
			break Loop1
		case XLS_TYPE_CODEPAGE:
			p, err := i.ReadCodePage(pos)
			if err != nil {
				return sheets, err
			}
			pos = p
		default:
			p, err := i.ReadDefault(pos)
			if err != nil {
				return sheets, err
			}
			pos = p
		}
	}

	fmt.Printf("fetch worksheet (%d) info.\n", len(i.Sheets))

	dataSize := len(i.Data)
	for index, sheet := range i.Sheets {
		if sheet.Type != 0x00 {
			// 0x00: Worksheet
			// 0x02: Chart
			// 0x06: Visual Basic module
			fmt.Printf("不是 Worksheet(%s): %d %d", sheet.Name, index, sheet.Type)
			continue
		}

		worksheet := &XlsBookWorkSheet{
			Name:             sheet.Name,
			LastColumnLetter: "A",
			LastColumnIdnex:  0,
			TotalRows:        0,
			TotalColumns:     0,
		}
		pos = int(sheet.Offset)
	Loop:
		for pos <= (dataSize - 4) {
			code, err := GetUInt2d(i.Data, pos)
			if err != nil {
				return nil, err
			}

			fmt.Printf("ListWookSheetNames Loop(1) at:%d  code: %d\n", pos, code)

			switch code {
			case XLS_TYPE_RK,
				XLS_TYPE_LABELSST,
				XLS_TYPE_NUMBER,
				XLS_TYPE_FORMULA,
				XLS_TYPE_BOOLERR,
				XLS_TYPE_LABEL:
				len, err := GetUInt2d(i.Data, pos+2)
				if err != nil {
					return nil, err
				}
				recordData, err := i.ReadRecordData(i.Data, pos+4, int(len))
				if err != nil {
					return nil, err
				}
				pos += 4 + int(len)
				// fmt.Printf("ListWookSheetNames Loop(2) at:%d size: %d\n", pos, len)
				rowIndex, err := GetUInt2d(recordData, 0)
				if err != nil {
					return nil, err
				}
				columnIndex, err := GetUInt2d(recordData, 2)
				if err != nil {
					return nil, err
				}
				worksheet.TotalRows = uint16(math.Max(float64(worksheet.TotalRows), float64(rowIndex+1)))
				worksheet.LastColumnIdnex = uint16(math.Max(float64(worksheet.LastColumnIdnex), float64(columnIndex)))
			case XLS_TYPE_BOF:
				p, err := i.ReadBof(pos)
				if err != nil {
					return sheets, err
				}
				pos = p
			case XLS_TYPE_EOF:
				fmt.Printf("ListWookSheetNames XLS_TYPE_EOF %d\n", pos)
				p, err := i.ReadDefault(pos)
				if err != nil {
					return sheets, err
				}
				pos = p
				break Loop
			default:
				p, err := i.ReadDefault(pos)
				if err != nil {
					return sheets, err
				}
				pos = p
			}
		}
		worksheet.LastColumnLetter = ColumnIndexToString(worksheet.LastColumnIdnex + 1)
		worksheet.TotalColumns = worksheet.LastColumnIdnex + 1
		sheets = append(sheets, worksheet)
	}
	return sheets, nil
}

func (i *XlsBook) ReadBof(pos int) (int, error) {
	length, err := GetUInt2d(i.Data, pos+2)
	if err != nil {
		return pos, err
	}

	fmt.Printf("ReadBof: %d %d\n", pos, length)

	recordData := i.Data[pos+4 : pos+4+int(length)]
	pos += 4 + int(length)
	substreamType, err := GetUInt2d(recordData, 2)
	if err != nil {
		return pos, err
	}

	fmt.Printf("ReadBof substreamType: %d\n", substreamType)

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
		fmt.Printf("xls 文件版本: %d\n", version)
		i.version = version
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
			p, err := i.ReadDefault(pos)
			if err != nil {
				return p, err
			}
			pos = p
			if code == XLS_TYPE_EOF || pos >= len(i.Data) {
				break
			}
		}
	}

	fmt.Printf("ReadBof end. %d\n", pos)
	return pos, nil
}

func (i *XlsBook) ReadSheet(pos int) (int, error) {
	length, err := GetUInt2d(i.Data, pos+2)
	if err != nil {
		return pos, err
	}
	recordData, err := i.ReadRecordData(i.Data, pos+4, int(length))
	if err != nil {
		return pos, err
	}

	recOffset, err := GetInt4d(i.Data, pos+4)
	if err != nil {
		return pos, err
	}
	fmt.Printf("sheet recOffset: %d\n", recOffset)
	pos += 4 + int(length)

	var sheetState string
	switch recordData[4] {
	case 0x00:
		sheetState = WORKSHEET_SHEETSTATE_VISIBLE
	case 0x01:
		sheetState = WORKSHEET_SHEETSTATE_HIDDEN
	case 0x02:
		sheetState = WORKSHEET_SHEETSTATE_VERYHIDDEN
	default:
		sheetState = WORKSHEET_SHEETSTATE_VISIBLE
	}
	fmt.Printf("sheet state: %s\n", sheetState)

	sheetType := recordData[5]
	fmt.Printf("sheet type: %d\n", sheetType)

	recName := ""
	switch i.version {
	case XLS_BIFF8:
		n, _, err := ReadUnicodeStringShort(recordData[6:])
		if err != nil {
			return pos, err
		}
		recName = string(n)
	case XLS_BIFF7:
		s, err := ReadByteStringStort(i.codePageCode, recordData[6:])
		if err != nil {
			return pos, err
		}
		recName = s
	}

	if i.Sheets == nil {
		i.Sheets = make([]XlsBookSheet, 0)
	}
	i.Sheets = append(i.Sheets, XlsBookSheet{
		Name:   recName,
		Offset: recOffset,
		State:  sheetState,
		Type:   sheetType,
	})
	return pos, nil
}

func (i *XlsBook) ReadDefault(pos int) (int, error) {
	len, err := GetUInt2d(i.Data, pos+2)
	if err != nil {
		return pos, err
	}
	return pos + 4 + int(len), nil
}

// func (i *XlsBook) ReadFilepass(pos int) {
// 	len := GetUInt2d(i.Data, pos+2)
// }

func (i *XlsBook) ReadCodePage(pos int) (int, error) {
	length, err := GetUInt2d(i.Data, pos+2)
	if err != nil {
		return pos, err
	}
	recordData, err := i.ReadRecordData(i.Data, pos+4, int(length))
	if err != nil {
		return pos, err
	}
	pos += 4 + int(length)
	codePage, err := GetUInt2d(recordData, 0)
	if err != nil {
		return pos, err
	}

	i.codePageCode = codePage
	i.codePage = NumberToName(codePage)

	return pos, nil
}

func (i *XlsBook) ReadRecordData(data []byte, pos int, len int) ([]byte, error) {
	data = data[pos : pos+len]

	// 没有加密，返回
	if i.encryption == MS_BIFF_CRYPTO_NONE || pos < int(i.encryptionStartPos) {
		fmt.Println("没加密")
		return data, nil
	}

	// XOR 不支持
	if i.encryption == MS_BIFF_CRYPTO_XOR {
		return nil, fmt.Errorf("不支持 XOR 加密")
	}

	// if i.encryption == MS_BIFF_CRYPTO_RC4 {
	// 	fmt.Println("RC4 加密")
	// 	oldBlock := int(math.Floor(float64(i.rc4Pos) / REKEY_BLOCK))
	// 	block := int(math.Floor(float64(pos) / REKEY_BLOCK))
	// 	endBlock := int(math.Floor(float64(pos + len) / REKEY_BLOCK))

	// TODO 这部分 PHPSpreadSheet 不完善
	// 	if block != oldBlock || pos < i.rc4Pos || i.rc4Cipher == nil {
	// 		i.rc4Cipher =
	// 	}
	// }

	return nil, fmt.Errorf("不明确的加密类型 %d", i.encryption)
}

func (i *XlsBook) MakeKey(block int, vc string) {

}
