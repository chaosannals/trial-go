package xlsrd4

type xlsColor string

const (
	// Colors
	COLOR_BLACK      xlsColor = "FF000000"
	COLOR_WHITE      xlsColor = "FFFFFFFF"
	COLOR_RED        xlsColor = "FFFF0000"
	COLOR_DARKRED    xlsColor = "FF800000"
	COLOR_BLUE       xlsColor = "FF0000FF"
	COLOR_DARKBLUE   xlsColor = "FF000080"
	COLOR_GREEN      xlsColor = "FF00FF00"
	COLOR_DARKGREEN  xlsColor = "FF008000"
	COLOR_YELLOW     xlsColor = "FFFFFF00"
	COLOR_DARKYELLOW xlsColor = "FF808000"
	COLOR_MAGENTA    xlsColor = "FFFF00FF"
	COLOR_CYAN       xlsColor = "FF00FFFF"
)

var (
	NAMED_COLORS = []string{
		"Black",
		"White",
		"Red",
		"Green",
		"Blue",
		"Yellow",
		"Magenta",
		"Cyan",
	}

	NAMED_COLOR_TRANSLATIONS = map[string]xlsColor{
		"Black":   COLOR_BLACK,
		"White":   COLOR_WHITE,
		"Red":     COLOR_RED,
		"Green":   COLOR_GREEN,
		"Blue":    COLOR_BLUE,
		"Yellow":  COLOR_YELLOW,
		"Magenta": COLOR_MAGENTA,
		"Cyan":    COLOR_CYAN,
	}

	INDEXED_COLORS = map[int]xlsColor{
		1:  "FF000000", //  System Colour #1 - Black
		2:  "FFFFFFFF", //  System Colour #2 - White
		3:  "FFFF0000", //  System Colour #3 - Red
		4:  "FF00FF00", //  System Colour #4 - Green
		5:  "FF0000FF", //  System Colour #5 - Blue
		6:  "FFFFFF00", //  System Colour #6 - Yellow
		7:  "FFFF00FF", //  System Colour #7- Magenta
		8:  "FF00FFFF", //  System Colour #8- Cyan
		9:  "FF800000", //  Standard Colour #9
		10: "FF008000", //  Standard Colour #10
		11: "FF000080", //  Standard Colour #11
		12: "FF808000", //  Standard Colour #12
		13: "FF800080", //  Standard Colour #13
		14: "FF008080", //  Standard Colour #14
		15: "FFC0C0C0", //  Standard Colour #15
		16: "FF808080", //  Standard Colour #16
		17: "FF9999FF", //  Chart Fill Colour #17
		18: "FF993366", //  Chart Fill Colour #18
		19: "FFFFFFCC", //  Chart Fill Colour #19
		20: "FFCCFFFF", //  Chart Fill Colour #20
		21: "FF660066", //  Chart Fill Colour #21
		22: "FFFF8080", //  Chart Fill Colour #22
		23: "FF0066CC", //  Chart Fill Colour #23
		24: "FFCCCCFF", //  Chart Fill Colour #24
		25: "FF000080", //  Chart Line Colour #25
		26: "FFFF00FF", //  Chart Line Colour #26
		27: "FFFFFF00", //  Chart Line Colour #27
		28: "FF00FFFF", //  Chart Line Colour #28
		29: "FF800080", //  Chart Line Colour #29
		30: "FF800000", //  Chart Line Colour #30
		31: "FF008080", //  Chart Line Colour #31
		32: "FF0000FF", //  Chart Line Colour #32
		33: "FF00CCFF", //  Standard Colour #33
		34: "FFCCFFFF", //  Standard Colour #34
		35: "FFCCFFCC", //  Standard Colour #35
		36: "FFFFFF99", //  Standard Colour #36
		37: "FF99CCFF", //  Standard Colour #37
		38: "FFFF99CC", //  Standard Colour #38
		39: "FFCC99FF", //  Standard Colour #39
		40: "FFFFCC99", //  Standard Colour #40
		41: "FF3366FF", //  Standard Colour #41
		42: "FF33CCCC", //  Standard Colour #42
		43: "FF99CC00", //  Standard Colour #43
		44: "FFFFCC00", //  Standard Colour #44
		45: "FFFF9900", //  Standard Colour #45
		46: "FFFF6600", //  Standard Colour #46
		47: "FF666699", //  Standard Colour #47
		48: "FF969696", //  Standard Colour #48
		49: "FF003366", //  Standard Colour #49
		50: "FF339966", //  Standard Colour #50
		51: "FF003300", //  Standard Colour #51
		52: "FF333300", //  Standard Colour #52
		53: "FF993300", //  Standard Colour #53
		54: "FF993366", //  Standard Colour #54
		55: "FF333399", //  Standard Colour #55
		56: "FF333333", //  Standard Colour #56
	}
)
