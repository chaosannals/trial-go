package xlsrd4

import (
	"fmt"
	"time"
)

type XlsSpreadsheetProperties struct {
	CodePage    uint16
	Title       string
	Subject     string
	Author      string
	Keywords    string
	Description string
	// Template 不支持
	LastModified string
	// Revision  不支持
	// Total Editing Time 不支持
	// Last Printed 不支持
	// Created Date/Time
	CreatedAt  time.Time
	ModifiedAt time.Time
	//    Number of Pages
	//    Number of Words
	//    Number of Characters
	//    Thumbnail
	//    Name of creating application
	//    Security

	Category string
	Manager  string
	Company  string
}

type XlsSpreadsheet struct {
	Properties XlsSpreadsheetProperties
}

type XlsSpreadsheetProperty struct {
	UInt16Value   uint16
	Int32Value    int32
	StringValue   string
	BoolValue     bool
	DateTimeValue time.Time
}

func (spreadsheet *XlsSpreadsheet) readSummaryInfo(book *XlsBook) error {
	//
	// offset: 0; size: 2; must be 0xFE 0xFF (UTF-16 LE byte order mark) UTF16的BOM
	// offset: 2; size: 2;
	// offset: 4; size: 2; OS version
	// offset: 6; size: 2; OS indicator
	// offset: 8; size: 16
	// offset: 24; size: 4; section count
	secCount, err := readInt4(book.summaryInfo, 24)
	if err != nil {
		return err
	}
	fmt.Printf("概要信息，片段数: %d\n", secCount)

	// offset: 28; size: 16; first section's class id: e0 85 9f f2 f9 4f 68 10 ab 91 08 00 2b 27 b3 d9
	// offset: 44; size: 4
	secOffset, err := readInt4(book.summaryInfo, 44)
	if err != nil {
		return err
	}
	fmt.Printf("概要信息，片起始位置: %d\n", secOffset)

	// section header
	// offset: $secOffset; size: 4; section length
	secLength, err := readInt4(book.summaryInfo, secOffset)
	if err != nil {
		return err
	}
	fmt.Printf("概要信息，片长度: %d\n", secLength)

	// offset: $secOffset+4; size: 4; property count
	propertiesCount, err := readInt4(book.summaryInfo, secOffset+4)
	if err != nil {
		return err
	}
	fmt.Printf("概要信息，属性个数: %d\n", propertiesCount)

	// initialize code page (used to resolve string values)
	codePage := uint16(1252) // 初始化 CP1252

	// offset: ($secOffset+8); size: var
	for i := int32(0); i < propertiesCount; i += 1 {
		// offset: ($secOffset+8) + (8 * $i); size: 4; property ID
		id, err := readInt4(book.summaryInfo, secOffset+8+(8*i))
		if err != nil {
			return err
		}
		fmt.Printf("概要信息，属性ID: %d\n", id)

		// offset: ($secOffset+12) + (8 * $i); size: 4; offset from beginning of section (48)
		offset, err := readInt4(book.summaryInfo, secOffset+12+(8*i))
		if err != nil {
			return err
		}

		// 类型
		t, err := readInt4(book.summaryInfo, secOffset+offset)
		if err != nil {
			return err
		}
		fmt.Printf("概要信息，属性类型: %d\n", t)

		value, err := GetXlsSpreadsheetProperty(book.summaryInfo, secOffset, offset, codePage, t)
		if err != nil {
			return err
		}

		switch id {
		case 0x01:
			codePage = value.UInt16Value
		case 0x02:
			spreadsheet.Properties.Title = value.StringValue
		case 0x03:
			spreadsheet.Properties.Subject = value.StringValue
		case 0x04:
			spreadsheet.Properties.Author = value.StringValue
		case 0x05:
			spreadsheet.Properties.Keywords = value.StringValue
		case 0x06:
			spreadsheet.Properties.Description = value.StringValue
		case 0x07: //    Template
		case 0x08: //    Last Saved By (LastModifiedBy)
			spreadsheet.Properties.LastModified = value.StringValue
		case 0x09: //    Revision
		case 0x0A: //    Total Editing Time
		case 0x0B: //    Last Printed
		case 0x0C:
			spreadsheet.Properties.CreatedAt = value.DateTimeValue
		case 0x0D:
			spreadsheet.Properties.ModifiedAt = value.DateTimeValue
		case 0x0E: //    Number of Pages
		case 0x0F: //    Number of Words
		case 0x10: //    Number of Characters
		case 0x11: //    Thumbnail
		case 0x12: //    Name of creating application
		case 0x13: //    Security
		}
	}
	return nil
}

func (spreadsheet *XlsSpreadsheet) readDocumentSummaryInfo(book *XlsBook) error {
	//    offset: 0;    size: 2;    must be 0xFE 0xFF (UTF-16 LE byte order mark)
	//    offset: 2;    size: 2;
	//    offset: 4;    size: 2;    OS version
	//    offset: 6;    size: 2;    OS indicator
	//    offset: 8;    size: 16
	//    offset: 24;    size: 4;    section count
	secCount, err := readInt4(book.documentSummaryInfo, 24)
	if err != nil {
		return err
	}
	fmt.Printf("文档概要信息，片段数: %d\n", secCount)

	// offset: 28;    size: 16;    first section's class id: 02 d5 cd d5 9c 2e 1b 10 93 97 08 00 2b 2c f9 ae
	// offset: 44;    size: 4;    first section offset
	secOffset, err := readInt4(book.documentSummaryInfo, 44)
	if err != nil {
		return err
	}
	fmt.Printf("文档概要信息，片起始位置: %d\n", secOffset)

	// section header
	// offset: $secOffset; size: 4; section length
	secLength, err := readInt4(book.documentSummaryInfo, secOffset)
	if err != nil {
		return err
	}
	fmt.Printf("文档概要信息，片长度: %d\n", secLength)

	// offset: $secOffset+4; size: 4; property count
	propertiesCount, err := readInt4(book.documentSummaryInfo, secOffset+4)
	if err != nil {
		return err
	}
	fmt.Printf("文档概要信息，属性个数: %d\n", propertiesCount)

	// initialize code page (used to resolve string values)
	codePage := uint16(1252) // 初始化 CP1252

	// offset: ($secOffset+8); size: var
	for i := int32(0); i < propertiesCount; i += 1 {
		// offset: ($secOffset+8) + (8 * $i); size: 4; property ID
		id, err := readInt4(book.documentSummaryInfo, secOffset+8+(8*i))
		if err != nil {
			return err
		}
		fmt.Printf("文档概要信息，属性ID: %d\n", id)

		// offset: ($secOffset+12) + (8 * $i); size: 4; offset from beginning of section (48)
		offset, err := readInt4(book.documentSummaryInfo, secOffset+12+(8*i))
		if err != nil {
			return err
		}

		// 类型
		t, err := readInt4(book.documentSummaryInfo, secOffset+offset)
		if err != nil {
			return err
		}
		fmt.Printf("文档概要信息，属性类型: %d\n", t)

		value, err := GetXlsSpreadsheetProperty(book.documentSummaryInfo, secOffset, offset, codePage, t)
		if err != nil {
			return err
		}

		switch id {
		case 0x01:
			codePage = value.UInt16Value
		case 0x02:
			spreadsheet.Properties.Category = value.StringValue
		case 0x03: //    Presentation Target
		case 0x04: //    Bytes
		case 0x05: //    Lines
		case 0x06: //    Paragraphs
		case 0x07: //    Slides
		case 0x08: //    Notes
		case 0x09: //    Hidden Slides
		case 0x0A: //    MM Clips
		case 0x0B: //    Scale Crop
		case 0x0C: //    Heading Pairs
		case 0x0D: //    Titles of Parts
		case 0x0E: //    Manager
			spreadsheet.Properties.Manager = value.StringValue
		case 0x0F:
			spreadsheet.Properties.Company = value.StringValue
		case 0x10: //    Links up-to-date
		}
	}
	return nil
}

// summaryInfo 和 documentSummaryInfo 可能是使用 2种的不同类别集，但是从代码上看基本相同，这里复用了。
func GetXlsSpreadsheetProperty(summaryInfo []byte, secOffset int32, offset int32, codePage uint16, t int32) (*XlsSpreadsheetProperty, error) {
	switch t {
	case 0x02:
		v, err := readUInt2(summaryInfo, secOffset+4+offset)
		return &XlsSpreadsheetProperty{
			UInt16Value: v,
		}, err
	case 0x03:
		v, err := readInt4(summaryInfo, secOffset+4+offset)
		return &XlsSpreadsheetProperty{
			Int32Value: v,
		}, err
	case 0x0B:
		v, err := readUInt2(summaryInfo, secOffset+4+offset)
		return &XlsSpreadsheetProperty{
			BoolValue: v != 0,
		}, err
	case 0x13:
		// 4 byte unsigned integer
		// PHPSpreadsheet 2015 年注释说之后要修复，2024 还没有改动。
		return nil, fmt.Errorf("不支持的类型 0x13")
	case 0x1E: // 这种是 C 字符串，字符集从代码上看是默认的 CP1252, 会在解析 ID 时候切换
		byteLength, err := readInt4(summaryInfo, secOffset+4+offset)
		if err != nil {
			return nil, err
		}
		d := summaryInfo[secOffset+8+offset : secOffset+8+offset+byteLength]
		v, err := convToUtf8ByCode(codePage, d)
		fmt.Printf("字符串：%s\n", v)
		return &XlsSpreadsheetProperty{
			StringValue: v,
		}, err
	case 0x40: // OLE 64位时间戳
		d := summaryInfo[secOffset+8+offset : secOffset+8+offset+8]
		ts, err := readOle2LocalDate(d)
		ut := time.Unix(ts, 0)
		fmt.Printf("OLE 时间：%s\n", ut.Format("2006-01-02 15:04:05"))
		return &XlsSpreadsheetProperty{
			DateTimeValue: ut,
		}, err
	case 0x47: // 剪切板格式
		return nil, fmt.Errorf("不支持的类型 0x47")
	case 0x100C: // documentSummaryInfo id 是 12 的数据， 被跳过了。
		return &XlsSpreadsheetProperty{}, nil
	case 0x101E: // documentSummaryInfo id 是 13 （标题之类的）的数据， 被跳过了。
		return &XlsSpreadsheetProperty{}, nil
	default:
		return nil, fmt.Errorf("不支持的类型 %d", t)
	}
}
