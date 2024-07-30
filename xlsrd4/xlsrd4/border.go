package xlsrd4

import "fmt"

type xlsBorderStyle string

const (
	// Border style
	BORDER_NONE             xlsBorderStyle = "none"
	BORDER_DASHDOT          xlsBorderStyle = "dashDot"
	BORDER_DASHDOTDOT       xlsBorderStyle = "dashDotDot"
	BORDER_DASHED           xlsBorderStyle = "dashed"
	BORDER_DOTTED           xlsBorderStyle = "dotted"
	BORDER_DOUBLE           xlsBorderStyle = "double"
	BORDER_HAIR             xlsBorderStyle = "hair"
	BORDER_MEDIUM           xlsBorderStyle = "medium"
	BORDER_MEDIUMDASHDOT    xlsBorderStyle = "mediumDashDot"
	BORDER_MEDIUMDASHDOTDOT xlsBorderStyle = "mediumDashDotDot"
	BORDER_MEDIUMDASHED     xlsBorderStyle = "mediumDashed"
	BORDER_SLANTDASHDOT     xlsBorderStyle = "slantDashDot"
	BORDER_THICK            xlsBorderStyle = "thick"
	BORDER_THIN             xlsBorderStyle = "thin"
	BORDER_OMIT             xlsBorderStyle = "omit" // should be used only for Conditional

	// Diagonal directions
	DIAGONAL_NONE = 0
	DIAGONAL_UP   = 1
	DIAGONAL_DOWN = 2
	DIAGONAL_BOTH = 3
)

var (
	BORDER_STYLE_MAP = map[int32]xlsBorderStyle{
		0x00: BORDER_NONE,
		0x01: BORDER_THIN,
		0x02: BORDER_MEDIUM,
		0x03: BORDER_DASHED,
		0x04: BORDER_DOTTED,
		0x05: BORDER_THICK,
		0x06: BORDER_DOUBLE,
		0x07: BORDER_HAIR,
		0x08: BORDER_MEDIUMDASHED,
		0x09: BORDER_DASHDOT,
		0x0A: BORDER_MEDIUMDASHDOT,
		0x0B: BORDER_DASHDOTDOT,
		0x0C: BORDER_MEDIUMDASHDOTDOT,
		0x0D: BORDER_SLANTDASHDOT,
	}
)

type xlsBorder struct {
	Style      xlsBorderStyle
	Color      xlsColor
	ColorIndex int32
}

func (border *xlsBorder) setStyle(i int32) {
	if v, ok := BORDER_STYLE_MAP[i]; ok {
		border.Style = v
	} else {
		border.Style = BORDER_NONE
		fmt.Printf("setStyle 无效的边框样式 %d", i)
	}
}

type xlsBorders struct {
	Left              xlsBorder
	Right             xlsBorder
	Top               xlsBorder
	Bottom            xlsBorder
	Diagonal          xlsBorder
	DiagonalDirection int
	AllBorders        xlsBorder
	Outline           xlsBorder
	Inside            xlsBorder
	Vertical          xlsBorder
	Horizontal        xlsBorder
}
