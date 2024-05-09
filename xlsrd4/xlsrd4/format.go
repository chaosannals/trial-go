package xlsrd4

import "fmt"

//  [MS-OI29500: Microsoft Office Implementation Information for ISO/IEC-29500 Standard Compliance]
//  18.8.30. numFmt (Number Format)
//
//  The ECMA standard defines built-in format IDs
//      14: "mm-dd-yy"
//      22: "m/d/yy h:mm"
//      37: "#,##0 ;(#,##0)"
//      38: "#,##0 ;[Red](#,##0)"
//      39: "#,##0.00;(#,##0.00)"
//      40: "#,##0.00;[Red](#,##0.00)"
//      47: "mmss.0"
//      KOR fmt 55: "yyyy-mm-dd"
//  Excel defines built-in format IDs
//      14: "m/d/yyyy"
//      22: "m/d/yyyy h:mm"
//      37: "#,##0_);(#,##0)"
//      38: "#,##0_);[Red](#,##0)"
//      39: "#,##0.00_);(#,##0.00)"
//      40: "#,##0.00_);[Red](#,##0.00)"
//      47: "mm:ss.0"
//      KOR fmt 55: "yyyy/mm/dd"

var xlsBuildInFormats map[uint16]string

func buildInFormatCode(index uint16) string {
	if xlsBuildInFormats == nil {
		xlsBuildInFormats = make(map[uint16]string)
		// General
		xlsBuildInFormats[0] = FORMAT_GENERAL
		xlsBuildInFormats[1] = "0"
		xlsBuildInFormats[2] = "0.00"
		xlsBuildInFormats[3] = "#,##0"
		xlsBuildInFormats[4] = "#,##0.00"

		xlsBuildInFormats[9] = "0%"
		xlsBuildInFormats[10] = "0.00%"
		xlsBuildInFormats[11] = "0.00E+00"
		xlsBuildInFormats[12] = "# ?/?"
		xlsBuildInFormats[13] = "# ??/??"
		xlsBuildInFormats[14] = FORMAT_DATE_XLSX14_ACTUAL // Despite ECMA 'mm-dd-yy'
		xlsBuildInFormats[15] = FORMAT_DATE_XLSX15
		xlsBuildInFormats[16] = "d-mmm"
		xlsBuildInFormats[17] = "mmm-yy"
		xlsBuildInFormats[18] = "h:mm AM/PM"
		xlsBuildInFormats[19] = "h:mm:ss AM/PM"
		xlsBuildInFormats[20] = "h:mm"
		xlsBuildInFormats[21] = "h:mm:ss"
		xlsBuildInFormats[22] = FORMAT_DATE_XLSX22_ACTUAL // Despite ECMA 'm/d/yy h:mm'

		xlsBuildInFormats[37] = "#,##0_);(#,##0)"            //  Despite ECMA '#,##0 ;(#,##0)'
		xlsBuildInFormats[38] = "#,##0_);[Red](#,##0)"       //  Despite ECMA '#,##0 ;[Red](#,##0)'
		xlsBuildInFormats[39] = "#,##0.00_);(#,##0.00)"      //  Despite ECMA '#,##0.00;(#,##0.00)'
		xlsBuildInFormats[40] = "#,##0.00_);[Red](#,##0.00)" //  Despite ECMA '#,##0.00;[Red](#,##0.00)'

		xlsBuildInFormats[44] = "_(\"$\"* #,##0.00_);_(\"$\"* \\(#,##0.00\\);_(\"$\"* \"-\"??_);_(@_)"
		xlsBuildInFormats[45] = "mm:ss"
		xlsBuildInFormats[46] = "[h]:mm:ss"
		xlsBuildInFormats[47] = "mm:ss.0" //  Despite ECMA 'mmss.0'
		xlsBuildInFormats[48] = "##0.0E+0"
		xlsBuildInFormats[49] = "@"

		// CHT
		xlsBuildInFormats[27] = "[$-404]e/m/d"
		xlsBuildInFormats[30] = "m/d/yy"
		xlsBuildInFormats[36] = "[$-404]e/m/d"
		xlsBuildInFormats[50] = "[$-404]e/m/d"
		xlsBuildInFormats[57] = "[$-404]e/m/d"

		// THA
		xlsBuildInFormats[59] = "t0"
		xlsBuildInFormats[60] = "t0.00"
		xlsBuildInFormats[61] = "t#,##0"
		xlsBuildInFormats[62] = "t#,##0.00"
		xlsBuildInFormats[67] = "t0%"
		xlsBuildInFormats[68] = "t0.00%"
		xlsBuildInFormats[69] = "t# ?/?"
		xlsBuildInFormats[70] = "t# ??/??"

		// JPN
		xlsBuildInFormats[28] = "[$-411]ggge\"年\"m\"月\"d\"日\""
		xlsBuildInFormats[29] = "[$-411]ggge\"年\"m\"月\"d\"日\""
		xlsBuildInFormats[31] = "yyyy\"年\"m\"月\"d\"日\""
		xlsBuildInFormats[32] = "h\"時\"mm\"分\""
		xlsBuildInFormats[33] = "h\"時\"mm\"分\"ss\"秒\""
		xlsBuildInFormats[34] = "yyyy\"年\"m\"月\""
		xlsBuildInFormats[35] = "m\"月\"d\"日\""
		xlsBuildInFormats[51] = "[$-411]ggge\"年\"m\"月\"d\"日\""
		xlsBuildInFormats[52] = "yyyy\"年\"m\"月\""
		xlsBuildInFormats[53] = "m\"月\"d\"日\""
		xlsBuildInFormats[54] = "[$-411]ggge\"年\"m\"月\"d\"日\""
		xlsBuildInFormats[55] = "yyyy\"年\"m\"月\""
		xlsBuildInFormats[56] = "m\"月\"d\"日\""
		xlsBuildInFormats[58] = "[$-411]ggge\"年\"m\"月\"d\"日\""
	}

	if r, ok := xlsBuildInFormats[index]; ok {
		return r
	} else {
		fmt.Printf("不是已知内建格式化串：%d", index)
		return FORMAT_GENERAL
	}
}
