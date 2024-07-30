package xlsrd4

import "fmt"

type xlsFillType string

const (
	FILL_NONE                    xlsFillType = "none"
	FILL_SOLID                   xlsFillType = "solid"
	FILL_GRADIENT_LINEAR         xlsFillType = "linear"
	FILL_GRADIENT_PATH           xlsFillType = "path"
	FILL_PATTERN_DARKDOWN        xlsFillType = "darkDown"
	FILL_PATTERN_DARKGRAY        xlsFillType = "darkGray"
	FILL_PATTERN_DARKGRID        xlsFillType = "darkGrid"
	FILL_PATTERN_DARKHORIZONTAL  xlsFillType = "darkHorizontal"
	FILL_PATTERN_DARKTRELLIS     xlsFillType = "darkTrellis"
	FILL_PATTERN_DARKUP          xlsFillType = "darkUp"
	FILL_PATTERN_DARKVERTICAL    xlsFillType = "darkVertical"
	FILL_PATTERN_GRAY0625        xlsFillType = "gray0625"
	FILL_PATTERN_GRAY125         xlsFillType = "gray125"
	FILL_PATTERN_LIGHTDOWN       xlsFillType = "lightDown"
	FILL_PATTERN_LIGHTGRAY       xlsFillType = "lightGray"
	FILL_PATTERN_LIGHTGRID       xlsFillType = "lightGrid"
	FILL_PATTERN_LIGHTHORIZONTAL xlsFillType = "lightHorizontal"
	FILL_PATTERN_LIGHTTRELLIS    xlsFillType = "lightTrellis"
	FILL_PATTERN_LIGHTUP         xlsFillType = "lightUp"
	FILL_PATTERN_LIGHTVERTICAL   xlsFillType = "lightVertical"
	FILL_PATTERN_MEDIUMGRAY      xlsFillType = "mediumGray"
)

var (
	FILL_PATTERN_MAP = map[uint32]xlsFillType{
		0x00: FILL_NONE,
		0x01: FILL_SOLID,
		0x02: FILL_PATTERN_MEDIUMGRAY,
		0x03: FILL_PATTERN_DARKGRAY,
		0x04: FILL_PATTERN_LIGHTGRAY,
		0x05: FILL_PATTERN_DARKHORIZONTAL,
		0x06: FILL_PATTERN_DARKVERTICAL,
		0x07: FILL_PATTERN_DARKDOWN,
		0x08: FILL_PATTERN_DARKUP,
		0x09: FILL_PATTERN_DARKGRID,
		0x0A: FILL_PATTERN_DARKTRELLIS,
		0x0B: FILL_PATTERN_LIGHTHORIZONTAL,
		0x0C: FILL_PATTERN_LIGHTVERTICAL,
		0x0D: FILL_PATTERN_LIGHTDOWN,
		0x0E: FILL_PATTERN_LIGHTUP,
		0x0F: FILL_PATTERN_LIGHTGRID,
		0x10: FILL_PATTERN_LIGHTTRELLIS,
		0x11: FILL_PATTERN_GRAY125,
		0x12: FILL_PATTERN_GRAY0625,
	}
)

type xlsFill struct {
	StartColorIndex uint16
	EndColorIndex   uint16
	FillType        xlsFillType
	Rotation        float64
	StartColor      xlsColor
	EndColor        xlsColor
	ColorChanged    bool
}

func (fill *xlsFill) setType(i uint32) {
	if t, ok := FILL_PATTERN_MAP[i]; ok {
		fill.FillType = t
	} else {
		fill.FillType = FILL_NONE
		fmt.Printf("xlsFill.setType 无效的索引 %d\n", i)
	}
}
