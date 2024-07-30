package xlsrd4

import "fmt"

type xlsBookParser struct {
	pos                int32
	version            uint16
	codePage           uint16
	encryption         int32
	encryptionStartPos int32

	excelCalendar int32

	fonts   []xlsFontInfo
	formats map[uint16]string

	styles   map[uint16]xlsStyleInfo
	xfStyles map[uint16]xlsStyleInfo

	xfIndex int

	isReadDataOnly bool // [功能]只读数据不读样式
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

	fmt.Printf("parseBof substreamType: %d\n", substreamType)

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

func (parser *xlsBookParser) parseDateMode(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return err
	}
	parser.pos += 4 + int32(length)

	parser.excelCalendar = CALENDAR_WINDOWS_1900
	if recordData[0] == 1 {
		parser.excelCalendar = CALENDAR_MAC_1904
	}
	fmt.Printf("DATE MODE(%d) %v: %d\n", length, recordData, parser.excelCalendar)
	return nil
}

func (parser *xlsBookParser) parseFont(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return err
	}
	parser.pos += 4 + int32(length)

	if !parser.isReadDataOnly {
		// offset: 0; size: 2; height of the font (in twips = 1/20 of a point)
		size, err := readUInt2(recordData, 0)
		if err != nil {
			return err
		}
		fmt.Printf("字体大小: %d\n", size/20)

		// offset: 2; size: 2; option flags
		// bit: 0; mask 0x0001; bold (redundant in BIFF5-BIFF8)
		// bit: 1; mask 0x0002; italic
		// bit: 2; mask 0x0004; underlined (redundant in BIFF5-BIFF8)
		// bit: 3; mask 0x0008; strikethrough
		flags, err := readUInt2(recordData, 2)
		if err != nil {
			return err
		}
		isBold := flags & 0x0001
		isItalic := flags & 0x0002 >> 1
		isUnderlined := flags & 0x0004 >> 2
		isStrikethrough := flags & 0x0008 >> 3

		fmt.Printf("字体 加粗: %d 斜体：%d 底线：%d 删除线: %d\n", isBold, isItalic, isUnderlined, isStrikethrough)

		// offset: 4; size: 2; colour index
		colorIndex, err := readUInt2(recordData, 4)
		if err != nil {
			return err
		}
		fmt.Printf("字体 颜色索引 %d\n", colorIndex)

		// offset: 6; size: 2; font weight
		// 该值 700 时认定字体 加粗
		fontWeight, err := readUInt2(recordData, 6)
		if err != nil {
			return err
		}
		fmt.Printf("字体 粗细 %d\n", fontWeight)

		// offset: 8; size: 2; escapement type
		escapementType, err := readUInt2(recordData, 8)
		if err != nil {
			return err
		}
		fmt.Printf("字体 escapementType %d\n", escapementType)

		// offset: 10; size: 1; underline type
		underlineType, err := readUInt2(recordData, 10)
		if err != nil {
			return err
		}
		fmt.Printf("字体 underlineType %d\n", underlineType)

		// offset: 11; size: 1; font family
		// offset: 12; size: 1; character set
		// offset: 13; size: 1; not used
		// offset: 14; size: var; font name
		var name string
		if parser.version == XLS_BIFF8 {
			v, _, err := readUnicodeStringShort(recordData[14:])
			if err != nil {
				return err
			}
			name = string(v)
			fmt.Printf("字体 BIFF8 名：%s\n", name)
		} else {
			v, err := readByteStringStort(parser.codePage, recordData[14:])
			if err != nil {
				return err
			}
			name = v
			fmt.Printf("字体 名：%s\n", v)
		}

		parser.fonts = append(parser.fonts, xlsFontInfo{
			Size:            size,
			IsBold:          isBold == 1,
			IsItalic:        isItalic == 1,
			IsUnderlined:    isUnderlined == 1,
			IsStrikethrough: isStrikethrough == 1,
			ColorIndex:      colorIndex,
			Weight:          fontWeight,
			EscapementType:  escapementType,
			UnderlineType:   underlineType,
			Name:            name,
		})
	}
	return nil
}

func (parser *xlsBookParser) parseFormat(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return err
	}
	parser.pos += 4 + int32(length)

	if !parser.isReadDataOnly {
		indexCode, err := readUInt2(recordData, 0)
		if err != nil {
			return err
		}
		fmt.Printf("FORMAT indexCode：%d\n", indexCode)

		var formatString string
		if parser.version == XLS_BIFF8 {
			v, _, err := readUnicodeStringShort(recordData[2:])
			if err != nil {
				return err
			}
			formatString = string(v)
			fmt.Printf("FORMAT STRING (BIFF8) ：%s\n", formatString)
		} else {
			v, err := readByteStringStort(parser.codePage, recordData[2:])
			if err != nil {
				return err
			}
			formatString = v
			fmt.Printf("FORMAT STRING：%s\n", v)
		}

		if parser.formats == nil {
			parser.formats = make(map[uint16]string)
		}

		parser.formats[indexCode] = formatString
	}
	return nil
}

func (parser *xlsBookParser) parseXf(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return err
	}
	parser.pos += 4 + int32(length)

	objStyle := xlsStyleInfo{}

	if !parser.isReadDataOnly {
		// offset:  0; size: 2; Index to FONT record
		fontIndex, err := readUInt2(recordData, 0)
		if err != nil {
			return err
		}
		// this has to do with that index 4 is omitted in all BIFF versions for some strange reason
		// check the OpenOffice documentation of the FONT record
		if fontIndex >= 4 {
			fontIndex -= 1
		}

		if int(fontIndex) < len(parser.fonts) {
			fmt.Printf("XF 字体： %s\n", parser.fonts[fontIndex].Name)
			objStyle.Font = parser.fonts[fontIndex]
		} else {
			fmt.Printf("XF 字体不匹配： %d\n", fontIndex)
		}

		// offset:  2; size: 2; Index to FORMAT record
		formatIndex, err := readUInt2(recordData, 2)
		if err != nil {
			return err
		}
		var formatCode string
		if f, ok := parser.formats[formatIndex]; ok {
			fmt.Printf("XF 格式化串: %s\n", f)
			formatCode = f
		} else {
			formatCode = buildInFormatCode(formatIndex)
			fmt.Printf("XF 格式化串(内建): %s\n", formatCode)
		}
		objStyle.NumberFormatCode = formatCode

		// offset:  4; size: 2; XF type, cell protection, and parent style XF
		// bit 2-0; mask 0x0007; XF_TYPE_PROT
		xfTypeProt, err := readUInt2(recordData, 4)
		if err != nil {
			return err
		}
		// bit 0; mask 0x01; 1 = cell is locked
		isLocked := xfTypeProt & 0x01
		if isLocked {
			objStyle.Protection.Locked = PROTECTION_INHERIT
		} else {
			objStyle.Protection.Locked = PROTECTION_UNPROTECTED
		}

		// bit 1; mask 0x02; 1 = Formula is hidden
		isHidden := (xfTypeProt & 0x02) >> 1
		if isHidden {
			objStyle.Protection.Hidden = PROTECTION_PROTECTED
		} else {
			objStyle.Protection.Hidden = PROTECTION_UNPROTECTED
		}
		// bit 2; mask 0x04; 0 = Cell XF, 1 = Cell Style XF
		isCellStyleXf := (xfTypeProt & 0x04) >> 2

		fmt.Printf("锁定：%d 隐藏: %d XF：%d\n", isLocked, isHidden, isCellStyleXf)

		// offset:  6; size: 1; Alignment and text break
		// bit 2-0, mask 0x07; horizontal alignment 横向对齐
		horAlign := recordData[6] & 0x07
		objStyle.Align.cellHorizontal(horAlign)

		// bit 3, mask 0x08; wrap text  换行
		wrapText := (recordData[6] & 0x08) >> 3
		if wrapText > 0 {
			objStyle.Align.WrapText = true
		}

		// bit 6-4, mask 0x70; vertical alignment 纵向对齐
		vertAlign := (recordData[6] & 0x70) >> 3
		objStyle.Align.cellVertical(vertAlign)

		fmt.Printf("横向对齐：%d 换行: %d 纵向对齐: %d\n", horAlign, wrapText, vertAlign)

		if parser.version == XLS_BIFF8 {
			// offset:  7; size: 1; XF_ROTATION: Text rotation angle
			angle := recordData[7]
			rotation := int16(0) // TODO 这个计算有问题。
			if angle <= 90 {
				rotation = int16(angle)
			} else if angle <= 180 {
				rotation = int16(90 - int(angle))
			} else if angle == TEXTROTATION_STACK_EXCEL {
				rotation = int16(TEXTROTATION_STACK_PHPSPREADSHEET)
			}
			fmt.Printf("旋转角度 %d\n", rotation)
			objStyle.Align.TextRotation = rotation

			// offset:  8; size: 1; Indentation, shrink to cell size, and text direction
			// bit: 3-0; mask: 0x0F; indent level
			indent := recordData[8] & 0x0F
			objStyle.Align.Indent = indent

			// bit: 4; mask: 0x10; 1 = shrink content to fit into cell
			shrinkToFit := recordData[8] & 0x10 >> 4
			if shrinkToFit > 0 {
				objStyle.Align.ShrinkToFit = true
			}

			fmt.Printf("XF indent: %d shrinkToFit: %d\n", indent, shrinkToFit)

			// offset:  9; size: 1; Flags used for attribute groups

			// offset: 10; size: 4; Cell border lines and background area
			// bit: 3-0; mask: 0x0000000F; left style
			bordersStyle, err := readInt4(recordData, 10)
			if err != nil {
				return err
			}
			leftStyle := bordersStyle & 0x0000000F
			objStyle.Borders.Left.setStyle(leftStyle)

			// bit: 7-4; mask: 0x000000F0; right style
			rightStyle := (bordersStyle & 0x000000F0) >> 4
			objStyle.Borders.Right.setStyle(rightStyle)

			// bit: 11-8; mask: 0x00000F00; top style
			topStyle := (bordersStyle & 0x00000F00) >> 8
			objStyle.Borders.Top.setStyle(topStyle)

			// bit: 15-12; mask: 0x0000F000; bottom style
			bottomStyle := (bordersStyle & 0x0000F000) >> 12
			objStyle.Borders.Bottom.setStyle(bottomStyle)

			// bit: 22-16; mask: 0x007F0000; left color
			leftColor := (bordersStyle & 0x007F0000) >> 16
			objStyle.Borders.Left.ColorIndex = leftColor

			// bit: 29-23; mask: 0x3F800000; right color
			rightColor := (bordersStyle & 0x3F800000) >> 23
			objStyle.Borders.Right.ColorIndex = rightColor

			// bit: 30; mask: 0x40000000; 1 = diagonal line from top left to right bottom
			diagonalDown := (bordersStyle & 0x40000000) >> 30
			// bit: 31; mask: 0x800000; 1 = diagonal line from bottom left to top right
			diagonalUp := (uint32(bordersStyle) & 0x80000000) >> 31

			if diagonalUp > 0 {
				if diagonalDown > 0 {
					objStyle.Borders.DiagonalDirection = DIAGONAL_BOTH
				} else {
					objStyle.Borders.DiagonalDirection = DIAGONAL_UP
				}
			} else {
				if diagonalDown > 0 {
					objStyle.Borders.DiagonalDirection = DIAGONAL_DOWN
				} else {
					objStyle.Borders.DiagonalDirection = DIAGONAL_NONE
				}
			}

			fmt.Printf("边框样式： 左: %d 右：%d 上：%d 下：%d 颜色：左：%d 右： %d 斜线: 下: %d 上：%d\n", leftStyle, rightStyle, topStyle, bottomStyle, leftColor, rightColor, diagonalDown, diagonalUp)

			// offset: 14; size: 4;
			// bit: 6-0; mask: 0x0000007F; top color
			bordersStyle2, err := readInt4(recordData, 14)
			if err != nil {
				return err
			}
			topColor := bordersStyle2 & 0x0000007F
			objStyle.Borders.Top.ColorIndex = topColor

			// bit: 13-7; mask: 0x00003F80; bottom color
			bottomColor := (bordersStyle2 & 0x00003F80) >> 7
			objStyle.Borders.Bottom.ColorIndex = bottomColor

			// bit: 20-14; mask: 0x001FC000; diagonal color
			diagonalColor := (bordersStyle2 & 0x001FC000) >> 14
			objStyle.Borders.Diagonal.ColorIndex = diagonalColor

			// bit: 24-21; mask: 0x01E00000; diagonal style
			diagonalStyle := (bordersStyle2 & 0x01E00000) >> 21
			objStyle.Borders.Diagonal.setStyle(diagonalStyle)

			// bit: 31-26; mask: 0xFC000000 fill pattern
			fillPattern := (uint32(bordersStyle2) & 0xFC000000) >> 28
			objStyle.Fill.setType(fillPattern)

			fmt.Printf("边框样式2： 颜色：上： %d 下 %d ：斜线 颜色 %d 样式: %d 填充模式：%d\n", topColor, bottomColor, diagonalColor, diagonalStyle, fillPattern)

			// offset: 18; size: 2; pattern and background colour
			// bit: 6-0; mask: 0x007F; color index for pattern color
			colorIndex, err := readUInt2(recordData, 18)
			if err != nil {
				return err
			}
			startColorIndex := colorIndex & 0x007F
			objStyle.Fill.StartColorIndex = startColorIndex

			// bit: 13-7; mask: 0x3F80; color index for pattern background
			endColorIndex := (colorIndex & 0x3F80) >> 7
			objStyle.Fill.EndColorIndex = endColorIndex

			fmt.Printf("前景色: %d 背景色: %d", startColorIndex, endColorIndex)
		} else {
			// BIFF5

			// offset: 7; size: 1; Text orientation and flags
			orientationAndFlags := recordData[7]

			// bit: 1-0; mask: 0x03; XF_ORIENTATION: Text orientation
			xfOrientation := 0x03 & orientationAndFlags
			rotation := 0
			switch xfOrientation {
			case 0:
				rotation = 0
			case 1:
				rotation = TEXTROTATION_STACK_PHPSPREADSHEET
			case 2:
				rotation = 90
			case 3:
				rotation = -90
			}
			objStyle.Align.TextRotation = rotation
			fmt.Printf("BIFF5 角度：%d\n", rotation)

			// offset: 8; size: 4; cell border lines and background area
			borderAndBackground, err := readInt4(recordData, 8)
			if err != nil {
				return err
			}

			// bit: 6-0; mask: 0x0000007F; color index for pattern color
			startColorIndex := borderAndBackground & 0x0000007F
			objStyle.Fill.StartColorIndex = uint16(startColorIndex)

			// bit: 13-7; mask: 0x00003F80; color index for pattern background
			endColorIndex := (borderAndBackground & 0x00003F80) >> 7
			objStyle.Fill.EndColorIndex = uint16(endColorIndex)

			// bit: 21-16; mask: 0x003F0000; fill pattern
			fillPatter := (borderAndBackground & 0x003F0000) >> 16
			objStyle.Fill.setType(uint32(fillPatter))

			// bit: 24-22; mask: 0x01C00000; bottom line style
			bottomLineStyle := (borderAndBackground & 0x01C00000) >> 22
			objStyle.Borders.Bottom.setStyle(bottomLineStyle)

			// bit: 31-25; mask: 0xFE000000; bottom line color
			bottomColorIndex := (uint32(borderAndBackground) & 0xFE000000) >> 25
			objStyle.Borders.Bottom.ColorIndex = int32(bottomColorIndex)

			fmt.Printf("BIFF5 前景色: %d 背景色： %d 填充色 %d 底部线样式 %d 颜色 %d\n", startColorIndex, endColorIndex, fillPatter, bottomLineStyle, bottomColorIndex)

			// TODO
			// offset: 12; size: 4; cell border lines
			borderLines, err := readInt4(recordData, 12)
			if err != nil {
				return err
			}
			// bit: 2-0; mask: 0x00000007; top line style
			topLineStyle := borderLines & 0x00000007
			// bit: 5-3; mask: 0x00000038; left line style
			leftLineStyle := (borderLines & 0x00000038) >> 3
			// bit: 8-6; mask: 0x000001C0; right line style
			rightLineStyle := (borderLines & 0x000001C0) >> 6
			// bit: 15-9; mask: 0x0000FE00; top line color index
			topLineColor := (borderLines & 0x0000FE00) >> 9
			// bit: 22-16; mask: 0x007F0000; left line color index
			leftLineColor := (borderLines & 0x007F0000) >> 16
			// bit: 29-23; mask: 0x3F800000; right line color index
			rightLineColor := (borderLines & 0x3F800000) >> 23

			fmt.Printf("线 样式： 上： %d 左 %d 右 %d  颜色 上 %d 左 %d 右 %d\n", topLineStyle, leftLineStyle, rightLineStyle, topLineColor, leftLineColor, rightLineColor)
		}

		// TODO 保存信息
		// if isCellStyleXf > 0 {
		// 	if parser.xfIndex == 0 {

		// 	}
		// } else {

		// }

		parser.xfIndex += 1
	}

	return nil
}

func (parser *xlsBookParser) parseXfExt(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return err
	}
	parser.pos += 4 + int32(length)

	if !parser.isReadDataOnly {
		// offset: 0; size: 2; 0x087D = repeated header
		// offset: 2; size: 2
		// offset: 4; size: 8; not used
		// offset: 12; size: 2; record version
		// offset: 14; size: 2; index to XF record which this record modifies
		ixfe, err := readUInt2(recordData, 14)
		if err != nil {
			return err
		}

		// offset: 16; size: 2; not used
		// offset: 18; size: 2; number of extension properties that follow
		cexts, err := readUInt2(recordData, 18)
		if err != nil {
			return err
		}
		// start reading the actual extension data
		offset := int32(20)
		for offset < int32(length) {
			// extension type
			extType, err := readUInt2(recordData, offset)
			if err != nil {
				return err
			}

			// extension length
			cb, err := readUInt2(recordData, offset+2)
			if err != nil {
				return err
			}

			// extension data
			extData := recordData[offset+4 : offset+4+int32(cb)]

			switch extType {
			case 4: // fill start color
				xclfType, err := readUInt2(extData, 0)
				if err != nil {
					return err
				}
				xclrValue := extData[4 : 4+4]

				if xclfType == 2 {
					rgb := fmt.Sprintf("%02X%02X%02X", xclrValue[0], xclrValue[1], xclrValue[2])
					fmt.Printf("RGB: %s\n", rgb)

					// modify the relevant style property
					// TODO
				}
			}
		}
	}

	return nil
}

func (parser *xlsBookParser) parseFilepass(workbook []byte) error {
	length, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	if length != 54 {
		return fmt.Errorf("File Pass 长度不是")
	}
	recordData, err := parser.parseRecordData(workbook, parser.pos+4, int32(length))
	if err != nil {
		return err
	}

	parser.pos += 4 + int32(length)

	// 验证密码 这个大概率没有被使用。
	err = verifyPassword([]byte("VelvetSweatshop"), recordData[6:22], recordData[22:38], recordData[38:54], []byte("md5Ctxt"))

	if err != nil {
		return err
	}

	parser.encryption = MS_BIFF_CRYPTO_RC4
	offset, err := readUInt2(workbook, parser.pos+2)
	if err != nil {
		return err
	}
	parser.encryptionStartPos = parser.pos + int32(offset)
	return nil
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

func verifyPassword(password []byte, docid []byte, saltData []byte, hashedSaltData []byte, valContext []byte) error {
	pwarray := make([]byte, 64)
	iMax := len(password)
	for i := 0; i < iMax; i += 1 {
		o := password[i]
		pwarray[2*i] = o
		pwarray[2*i+1] = 0
	}

	pwarray[2*(iMax-1)] = 0x80
	pwarray[56] = ((byte(iMax) - 1) << 4) & 0xFF

	// TODO
	return nil
}
