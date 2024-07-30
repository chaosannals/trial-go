package xlsrd4

import "fmt"

type xlsAlignment string

const (
	HORIZONTAL_GENERAL           xlsAlignment = "general"
	HORIZONTAL_LEFT              xlsAlignment = "left"
	HORIZONTAL_RIGHT             xlsAlignment = "right"
	HORIZONTAL_CENTER            xlsAlignment = "center"
	HORIZONTAL_CENTER_CONTINUOUS xlsAlignment = "centerContinuous"
	HORIZONTAL_JUSTIFY           xlsAlignment = "justify"
	HORIZONTAL_FILL              xlsAlignment = "fill"
	HORIZONTAL_DISTRIBUTED       xlsAlignment = "distributed" // Excel2007 only
)

var (
	// Mapping for horizontal alignment
	HORIZONTAL_ALIGNMENT_FOR_XLSX = map[xlsAlignment]xlsAlignment{
		HORIZONTAL_LEFT:              HORIZONTAL_LEFT,
		HORIZONTAL_RIGHT:             HORIZONTAL_RIGHT,
		HORIZONTAL_CENTER:            HORIZONTAL_CENTER,
		HORIZONTAL_CENTER_CONTINUOUS: HORIZONTAL_CENTER_CONTINUOUS,
		HORIZONTAL_JUSTIFY:           HORIZONTAL_JUSTIFY,
		HORIZONTAL_FILL:              HORIZONTAL_FILL,
		HORIZONTAL_DISTRIBUTED:       HORIZONTAL_DISTRIBUTED,
	}

	// Mapping for horizontal alignment CSS
	HORIZONTAL_ALIGNMENT_FOR_HTML = map[xlsAlignment]xlsAlignment{
		HORIZONTAL_LEFT:              HORIZONTAL_LEFT,
		HORIZONTAL_RIGHT:             HORIZONTAL_RIGHT,
		HORIZONTAL_CENTER:            HORIZONTAL_CENTER,
		HORIZONTAL_CENTER_CONTINUOUS: HORIZONTAL_CENTER,
		HORIZONTAL_JUSTIFY:           HORIZONTAL_JUSTIFY,
		//HORIZONTAL_FILL : HORIZONTAL_FILL, // no reasonable equivalent for fill
		HORIZONTAL_DISTRIBUTED: HORIZONTAL_JUSTIFY,
	}
)

const (
	VERTICAL_BOTTOM      = "bottom"
	VERTICAL_TOP         = "top"
	VERTICAL_CENTER      = "center"
	VERTICAL_JUSTIFY     = "justify"
	VERTICAL_DISTRIBUTED = "distributed" // Excel2007 only
	// Vertical alignment CSS
	VERTICAL_BASELINE    = "baseline"
	VERTICAL_MIDDLE      = "middle"
	VERTICAL_SUB         = "sub"
	VERTICAL_SUPER       = "super"
	VERTICAL_TEXT_BOTTOM = "text-bottom"
	VERTICAL_TEXT_TOP    = "text-top"
)

var (
	// Mapping for vertical alignment
	VERTICAL_ALIGNMENT_FOR_XLSX = map[xlsAlignment]xlsAlignment{
		VERTICAL_BOTTOM:      VERTICAL_BOTTOM,
		VERTICAL_TOP:         VERTICAL_TOP,
		VERTICAL_CENTER:      VERTICAL_CENTER,
		VERTICAL_JUSTIFY:     VERTICAL_JUSTIFY,
		VERTICAL_DISTRIBUTED: VERTICAL_DISTRIBUTED,
		// css settings that arent't in sync with Excel
		VERTICAL_BASELINE:    VERTICAL_BOTTOM,
		VERTICAL_MIDDLE:      VERTICAL_CENTER,
		VERTICAL_SUB:         VERTICAL_BOTTOM,
		VERTICAL_SUPER:       VERTICAL_TOP,
		VERTICAL_TEXT_BOTTOM: VERTICAL_BOTTOM,
		VERTICAL_TEXT_TOP:    VERTICAL_TOP,
	}
)

var (
	CELL_ALIGNMENT_HORIZONTAL = map[int]xlsAlignment{
		0: HORIZONTAL_GENERAL,
		1: HORIZONTAL_LEFT,
		2: HORIZONTAL_CENTER,
		3: HORIZONTAL_RIGHT,
		4: HORIZONTAL_FILL,
		5: HORIZONTAL_JUSTIFY,
		6: HORIZONTAL_CENTER_CONTINUOUS,
	}
	CELL_ALIGNMENT_VERTICAL = map[int]xlsAlignment{
		0: VERTICAL_TOP,
		1: VERTICAL_CENTER,
		2: VERTICAL_BOTTOM,
		3: VERTICAL_JUSTIFY,
	}
)

type xlsAlign struct {
	Horizontal   xlsAlignment
	Vertical     xlsAlignment
	TextRotation int
	WrapText     bool
	ShrinkToFit  bool
	Indent       int
	ReadOrder    int
}

func (align *xlsAlign) cellHorizontal(i int) {
	if h, ok := CELL_ALIGNMENT_HORIZONTAL[i]; ok {
		align.Horizontal = h
	} else {
		fmt.Printf("cellHorizontal 出现无效值 %d", i)
	}
}

func (align *xlsAlign) cellVertical(i int) {
	if v, ok := CELL_ALIGNMENT_VERTICAL[i]; ok {
		align.Vertical = v
	} else {
		fmt.Printf("cellVertical 出现无效值 %d", i)
	}
}
