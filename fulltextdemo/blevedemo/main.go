package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
	bi "github.com/blevesearch/bleve_index_api"
	"github.com/go-faker/faker/v4"
)

type FooDoc struct {
	Id   string `json:"id" faker:"uuid_digit"`
	From string `json:"from" faker:"chinese_name"`
	// Body string `json:"body" faker:"lang=chi,sentence"`
	Body string `json:"body" faker:"paragraph"`
	Age  int32  `json:"age" faker:"oneof: 4, 64"`
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getIndex() (bleve.Index, error) {
	indexDir := "foo.bleve"

	if IsFileExist(indexDir) {
		return bleve.Open(indexDir)
	} else {
		// 很麻烦，需要代码定义结构
		mapping := bleve.NewIndexMapping() // 索引映射
		// docMapping := bleve.NewDocumentMapping()
		// mapping.AddDocumentMapping("foo", docMapping)   // 文档映射
		// fieldBodyMapping := bleve.NewTextFieldMapping() // body 字段映射
		// fieldBodyMapping.Store = true
		// docMapping.AddFieldMappingsAt("body", fieldBodyMapping)
		return bleve.New(indexDir, mapping)
	}
}

// 间隔 1秒 添加 10000 条 faker 数据。
func addDoc(index bleve.Index) {
	// 不可以有多个 Index 被同时操作, 会锁死。
	// index, err := getIndex()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer index.Close()

	for {
		fmt.Printf("生成开始 %s\n", time.Now().Format("2006-01-02 15:04:05"))
		batch := index.NewBatch()
		var doc FooDoc
		for i := 0; i < 10000; i++ {
			if err := faker.FakeData(&doc); err != nil {
				log.Fatal(err)
			}
			if err := batch.Index(doc.Id, doc); err != nil {
				log.Fatal(err)
			}
		}

		if err := index.Batch(batch); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("生成结束 %s\n", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(1 * time.Second)
	}
}

// 每秒随机查一次
func searchDoc(index bleve.Index) {
	// 不可以有多个 Index 被同时操作, 会锁死。
	// index, err := getIndex()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer index.Close()

	fmt.Println("生成随机关键字：================================================")
	// 选项无效，大概没实现。
	// p := faker.Sentence(func(oo *options.Options) {
	// 	oo.StringLanguage = &interfaces.LangCHI
	// })
	p := faker.ChineseName()

	fmt.Printf("%s\n", p)

	query := bleve.NewMatchQuery(p)
	search := bleve.NewSearchRequest(query)
	search.Explain = false
	search.Size = 1000
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}

	hitCount := searchResults.Hits.Len()
	fmt.Printf("result: %v\n", searchResults)
	fmt.Printf("facets: %v\n", searchResults.Facets)
	for i := 0; i < hitCount; i++ {
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

		// 这里没有啥有用的信息。
		fmt.Printf("fields: %d %d %d %s %s\n", i, len(hit.Fields), hit.Size(), hit.ID, hit.Index)
		for k, v := range hit.Fields {
			fmt.Printf("[%d] %s => %v\n", i, k, v)
		}
		fmt.Printf("hits: %d\n", hitCount)
	}

	fmt.Println("搜索结束：================================================")
}

func main() {
	// 只要是同个 index ，被多个协程使用是容许的。
	index, err := getIndex()
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	go addDoc(index)

	for {
		fmt.Println("搜索：")
		searchDoc(index)
		time.Sleep(1 * time.Second)
	}
}
