package xlsrd4

import "fmt"

type xlsBookParser struct {
	pos                int32
	version            uint16
	codePage           uint16
	encryption         int32
	encryptionStartPos int32
}

func (parser *xlsBookParser) parseSheet(workbook []byte) (*xlsBookSheet, error) {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return nil, err
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return nil, err
	}

	recOffset, err := readInt4(workbook, parser.pos+4)
	if err != nil {
		return nil, err
	}
	fmt.Printf("sheet recOffset: %d\n", recOffset)
	parser.pos += 4 + int32(length)

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
	switch parser.version {
	case XLS_BIFF8:
		n, _, err := readUnicodeStringShort(recordData[6:])
		if err != nil {
			return nil, err
		}
		recName = string(n)
	case XLS_BIFF7:
		s, err := readByteStringStort(parser.codePage, recordData[6:])
		if err != nil {
			return nil, err
		}
		recName = s
	}

	return &xlsBookSheet{
		Name:   recName,
		Offset: recOffset,
		State:  sheetState,
		Type:   sheetType,
	}, nil
}

// 获取头部信息，某些版本可以得到一个版本号。
func (parser *xlsBookParser) parseBof(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}

	fmt.Printf("parseBof: %d %d\n", parser.pos, length)

	recordData := workbook[parser.pos+4 : parser.pos+4+int32(length)]
	parser.pos += 4 + int32(length)
	substreamType, err := readUInt2(recordData, 2)
	if err != nil {
		return err
	}

	fmt.Printf("ReadBof substreamType: %d\n", substreamType)

	switch substreamType {
	case XLS_WORKBOOKGLOBALS:
		parser.version, err = readUInt2(recordData, 0)
		if err != nil {
			return err
		}
		if (parser.version != XLS_BIFF8) && (parser.version != XLS_BIFF7) {
			return fmt.Errorf("这个文件的格式太老了，不支持。")
		}
		fmt.Printf("xls 文件版本: %d\n", parser.version)
		return nil
	case XLS_WORKSHEET:
		// 此项的版本信息不可靠
		// (OpenOffice doc, 5.8)指出使用全局流里面的版本信息
		fmt.Println("xls 文件此种类型非正规")
		return nil
	default:
		for {
			code, err := readUInt2(workbook, parser.pos)
			if err != nil {
				return err
			}
			err = parser.parseSkip(workbook)
			if err != nil {
				return err
			}
			if code == XLS_TYPE_EOF || parser.pos >= int32(len(workbook)) {
				break
			}
		}
	}

	fmt.Printf("ReadBof end. %d\n", parser.pos)
	return nil
}

func (parser *xlsBookParser) parseCodePage(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return err
	}
	parser.pos += 4 + int32(length)
	parser.codePage, err = readUInt2(recordData, 0)
	return err
}

func (parser *xlsBookParser) parseRecordData(workbook []byte, pos int32, length int32) ([]byte, error) {
	data := workbook[pos : pos+length]

	// 没有加密，返回
	if parser.encryption == MS_BIFF_CRYPTO_NONE || pos < parser.encryptionStartPos {
		// fmt.Println("没加密")
		return data, nil
	}
	fmt.Println("加密数据")

	// XOR 不支持
	if parser.encryption == MS_BIFF_CRYPTO_XOR {
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

	return nil, fmt.Errorf("不明确的加密类型 %d", parser.encryption)
}

// 在读取不同信息时用于跳过块的解析
func (parser *xlsBookParser) parseSkip(workbook []byte) error {
	len, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	parser.pos += 4 + int32(len)
	return nil
}
