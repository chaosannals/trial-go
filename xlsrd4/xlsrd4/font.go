package xlsrd4

type xlsFontInfo struct {
	Size            uint16
	IsBold          bool
	IsItalic        bool
	IsUnderlined    bool
	IsStrikethrough bool
	Color           xlsColor
	ColorIndex      uint16
	Weight          uint16
	EscapementType  uint16
	UnderlineType   uint16
	Name            string
}
