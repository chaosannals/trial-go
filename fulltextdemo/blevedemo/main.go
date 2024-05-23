package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	bi "github.com/blevesearch/bleve_index_api"
)

type FooDoc struct {
	Id   string `json:"id"`
	From string `json:"from"`
	Body string `json:"body"`
	Age  int32  `json:"age"`
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func main() {
	indexDir := "foo.bleve"

	var index bleve.Index
	var err error

	if IsFileExist(indexDir) {
		index, err = bleve.Open(indexDir)
	} else {
		// 很麻烦，需要代码定义结构
		mapping := bleve.NewIndexMapping() // 索引映射
		// docMapping := bleve.NewDocumentMapping()
		// mapping.AddDocumentMapping("foo", docMapping)   // 文档映射
		// fieldBodyMapping := bleve.NewTextFieldMapping() // body 字段映射
		// fieldBodyMapping.Store = true
		// docMapping.AddFieldMappingsAt("body", fieldBodyMapping)
		index, err = bleve.New(indexDir, mapping)
	}
	if err != nil {
		log.Fatal(err)
	}

	doc := FooDoc{
		Id:   "0002",
		From: "text中文，分词",
		Body: "测试",
		Age:  0xCDEF,
	}

	index.Index(doc.Id, doc)

	query := bleve.NewMatchQuery("中文")
	search := bleve.NewSearchRequest(query)
	search.Explain = false
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("result: %v\n", searchResults)
	fmt.Printf("facets: %v\n", searchResults.Facets)
	fmt.Printf("hits: %d\n", searchResults.Hits.Len())
	for i := 0; i < searchResults.Hits.Len(); i++ {
		hit := searchResults.Hits[i]
		// fmt.Printf("expl: %v", hit.Expl)
		if doc, err := index.Document(hit.ID); err != nil {
			fmt.Printf("err: %v", err)
		} else {
			fmt.Printf("doc: %v\n", doc)

			doc.VisitFields(func(f bi.Field) { // 要靠遍历器遍历。
				n := f.Name()
				// 值类型要靠自己区分,拿到的是存储块 []byte。
				switch t := f.(type) {
				case bi.NumericField: // 这个是 float64 没有整型，整型全被转成浮点。
					if v, err := t.Number(); err != nil {
						fmt.Printf("转换数字错误： %v\n", err)
					} else {
						fmt.Printf("数字: %s %f\n", n, v)
					}
				default:
					// 没看到文档明确标注大小头。本机 x86 无论大小头读，数据都不对。
					// 这个应该是浮点 float64。 有点恶心，浮点是机器相关，直接用浮点存文件。。。。。
					// 不能直接通过下面去获取整型值。
					// b := f.Value()
					// v := binary.LittleEndian.Uint64(b)
					// fmt.Printf("doc field %s %v %v \n", n, v, b)

					v := string(f.Value())
					fmt.Printf("doc field %s %v\n", n, v)
				}
			})
		}

		fmt.Printf("fields: %d %d %d %s %s\n", i, len(hit.Fields), hit.Size(), hit.ID, hit.Index)
		for k, v := range hit.Fields {
			fmt.Printf("[%d] %s => %v\n", i, k, v)
		}
	}
}
