package models

import (
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"os"
	"path/filepath"
	"strings"
)

// 获取字典路径
func getDictPath() string {
	result := make([]string, 2)
	root := filepath.Dir(os.Args[0])
	primary := filepath.Join(root, "dictionary.txt")
	if _, e := os.Stat(primary); e == nil || os.IsExist(e) {
		result = append(result, primary)
	}
	result = append(result, "zh")
	return strings.Join(result, ",")
}

//AdpotIndex 生成一个字段的索引
func AdpotIndex(name string) *riot.Engine {
	searcher := riot.Engine{}
	searcher.Init(types.EngineOpts{
		Using:       3,
		GseDict:     getDictPath(),
		UseStore:    true,
		StoreFolder: filepath.Join("data", name),
	})
	return &searcher
}
