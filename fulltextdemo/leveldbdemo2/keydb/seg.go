package keydb

import (
	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/pos"
)

var (
	Seg    gse.Segmenter
	PosSeg pos.Segmenter
)

func InitSeg() {
	Seg.LoadDict() // 加载默认字典
	// 加载默认 embed 词典
	// seg.LoadDictEmbed()
	//
	// 加载简体中文词典
	Seg.LoadDict("zh_s")
	Seg.LoadDictEmbed("zh_s")
}
