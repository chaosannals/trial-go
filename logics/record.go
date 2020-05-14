package logics

import (
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"os"
	"path/filepath"
	"strings"
)

var searcher = riot.Engine{}

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

func Init() func() {
	searcher.Init(types.EngineOpts{
		Using:       3,
		GseDict:     getDictPath(),
		UseStore:    true,
		StoreFolder: "data",
	})
	return func() {
		searcher.Close()
	}
}

func Change(id string, data types.DocData, force ...bool) {
	searcher.Index(id, data, force...)
	searcher.Flush()
}

func Remove(id string) {
	searcher.RemoveDoc(id, true)
	searcher.Flush()
}

func Search(request types.SearchReq) types.SearchResp {
	return searcher.Search(request)
}
