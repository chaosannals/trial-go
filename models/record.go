package models

import (
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var searcher = riot.Engine{}

func Init() func() {
	searcher.Init(types.EngineOpts{
		Using:   3,
		GseDict: "zh",
	})
	return func() {
		searcher.Close()
	}
}

func Change(id string, data types.DocData, force ...bool) {
	searcher.Index(id, data, force...)
	searcher.Flush()
}

func Search(request types.SearchReq) types.SearchResp {
	return searcher.Search(request)
}
