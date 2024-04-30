package xlsrd4

import (
	"fmt"
	"math"
)

type XlsSheetInfo struct {
	Name             string
	LastColumnLetter string
	LastColumnIndex  uint16
	TotalRows        uint16
	TotalColumns     uint16
}

type XlsBook struct {
	oleFile             *xlsOleFile
	workbook            []byte
	summaryInfo         []byte
	documentSummaryInfo []byte
}

type xlsBookSheet struct {
	Name   string
	Offset int32
	State  string
	Type   uint8
}

func ReadXlsFile(xlsPath string) (*XlsBook, error) {
	oleFile, err := readOleFile(xlsPath)
	if err != nil {
		return nil, err
	}

	workbook, err := oleFile.readStream(oleFile.propertySet.WorkBookId)
	if err != nil {
		return nil, err
	}

	summaryInfo, err := oleFile.readStream(oleFile.propertySet.SummaryInfoId)
	if err != nil {
		return nil, err
	}

	documentSummaryInfo, err := oleFile.readStream(oleFile.propertySet.DocumentSummaryInfoId)
	if err != nil {
		return nil, err
	}

	return &XlsBook{
		oleFile:             oleFile,
		workbook:            workbook,
		summaryInfo:         summaryInfo,
		documentSummaryInfo: documentSummaryInfo,
	}, nil
}

func ListXlsSheetInfo(xlsPath string) ([]XlsSheetInfo, error) {
	book, err := ReadXlsFile(xlsPath)
	if err != nil {
		return nil, err
	}
	dLen := int32(len(book.workbook))
	bookSheets := make([]xlsBookSheet, 0)
	parser := &xlsBookParser{
		pos: 0,
	}

	fmt.Printf("workbook size: %d\n", dLen)
	// 解析过程获取 BookSheet
parseLoop:
	for parser.pos < dLen {
		// fmt.Printf("ListXlsSheetInfo parseLoop: %d\n", parser.pos)
		code, err := readUInt2(book.workbook, parser.pos)
		if err != nil {
			return nil, err
		}
		// fmt.Printf("ListXlsSheetInfo code: %d %d\n", parser.pos, code)
		switch code {
		case XLS_TYPE_BOF:
			err = parser.parseBof(book.workbook)
			if err != nil {
				return nil, err
			}
		case XLS_TYPE_SHEET:
			sheet, err := parser.parseSheet(book.workbook)
			if err != nil {
				return nil, err
			}
			bookSheets = append(bookSheets, *sheet)
		case XLS_TYPE_EOF:
			err := parser.parseSkip(book.workbook)
			if err != nil {
				return nil, err
			}
			break parseLoop
		case XLS_TYPE_CODEPAGE:
			err = parser.parseCodePage(book.workbook)
			if err != nil {
				return nil, err
			}
		default:
			err := parser.parseSkip(book.workbook)
			if err != nil {
				return nil, err
			}
			// fmt.Printf("pos: %d size: %d", parser.pos, dLen)
		}
	}

	// 获取信息
	sheets := make([]XlsSheetInfo, 0)
	dataSize := int32(len(book.workbook))
	for index, sheet := range bookSheets {
		if sheet.Type != 0x00 {
			// 0x00: Worksheet
			// 0x02: Chart
			// 0x06: Visual Basic module
			fmt.Printf("不是 Worksheet(%s): %d %d", sheet.Name, index, sheet.Type)
			continue
		}

		worksheet := &XlsSheetInfo{
			Name:             sheet.Name,
			LastColumnLetter: "A",
			LastColumnIndex:  0,
			TotalRows:        0,
			TotalColumns:     0,
		}
		parser.pos = sheet.Offset
	Loop:
		for parser.pos <= (dataSize - 4) {
			code, err := readUInt2(book.workbook, parser.pos)
			if err != nil {
				return nil, err
			}

			// fmt.Printf("ListWookSheetNames Loop(1) at:%d  code: %d\n", pos, code)

			switch code {
			case XLS_TYPE_RK,
				XLS_TYPE_LABELSST,
				XLS_TYPE_NUMBER,
				XLS_TYPE_FORMULA,
				XLS_TYPE_BOOLERR,
				XLS_TYPE_LABEL:
				len, err := readUInt2(book.workbook, parser.pos+2)
				if err != nil {
					return nil, err
				}
				recordData, err := parser.parseRecordData(book.workbook, parser.pos+4, int32(len))
				if err != nil {
					return nil, err
				}
				parser.pos += 4 + int32(len)
				// fmt.Printf("ListWookSheetNames Loop(2) at:%d size: %d\n", pos, len)
				rowIndex, err := readUInt2(recordData, 0)
				if err != nil {
					return nil, err
				}
				columnIndex, err := readUInt2(recordData, 2)
				if err != nil {
					return nil, err
				}
				worksheet.TotalRows = uint16(math.Max(float64(worksheet.TotalRows), float64(rowIndex+1)))
				worksheet.LastColumnIndex = uint16(math.Max(float64(worksheet.LastColumnIndex), float64(columnIndex)))
			case XLS_TYPE_BOF:
				err := parser.parseBof(book.workbook)
				if err != nil {
					return sheets, err
				}
			case XLS_TYPE_EOF:
				fmt.Printf("ListWookSheetNames XLS_TYPE_EOF %d\n", parser.pos)
				err := parser.parseSkip(book.workbook)
				if err != nil {
					return sheets, err
				}
				break Loop
			default:
				err := parser.parseSkip(book.workbook)
				if err != nil {
					return sheets, err
				}
			}
		}
		worksheet.LastColumnLetter = ColumnIndexToString(worksheet.LastColumnIndex + 1)
		worksheet.TotalColumns = worksheet.LastColumnIndex + 1
		sheets = append(sheets, *worksheet)
	}

	return sheets, nil
}

func LoadSpreadsheetFromFile(xlsPath string) (*XlsSpreadsheet, error) {
	book, err := ReadXlsFile(xlsPath)
	if err != nil {
		return nil, err
	}
	spreadsheet := &XlsSpreadsheet{}
	err = spreadsheet.readSummaryInfo(book)
	return spreadsheet, err
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
