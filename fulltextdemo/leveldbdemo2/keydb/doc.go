package keydb

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	mapset "github.com/deckarep/golang-set/v2"
)

type DocContent struct {
	Content   string  `json:"content" form:"content"`
	IntInfo   int32   `json:"intInfo" form:"intInfo"`
	FloatInfo float64 `json:"floatInfo" form:"floatInfo"`
}
type DocStore struct {
	DocContent
	ID uuid.UUID `json:"id" form:"id"`
}

type DocIndex = []uuid.UUID

func (doc *DocContent) InsertAndCut() (uuid.UUID, []string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return id, nil, err
	}

	de, err := json.Marshal(doc)
	if err != nil {
		return id, nil, err
	}

	// 索引
	// seqs := gse.ToString(keydb.Seg.Segment([]byte(doc.Content)), true)
	seqs := Seg.CutSearch(doc.Content, true)
	for _, s := range seqs {
		k := []byte(s)
		r := id[:]
		if v, err := IndexDb.Get(k, &opt.ReadOptions{}); err != nil {
			if err == leveldb.ErrNotFound {
				fmt.Printf("插入 %v\n", r)
			} else {
				fmt.Printf("插入错误 %v\n", err)
			}
		} else {
			r = append(r, v...)
		}
		IndexDb.Put(k, r, &opt.WriteOptions{})
	}

	// 文档
	err = DocDb.Put(id[:], de, &opt.WriteOptions{})

	return id, seqs, err
}

// 切词 AND 逻辑
func Query(text string) ([]DocStore, error) {
	result := make([]DocStore, 0)

	fmt.Printf("检索： %s\n", text)
	seqs := Seg.CutSearch(text, true)

	// 检索索引
	ids := mapset.NewSet[uuid.UUID]()
	for _, seq := range seqs {
		k := []byte(seq)
		if v, err := IndexDb.Get(k, &opt.ReadOptions{}); err != nil {
			if err == leveldb.ErrNotFound {
				fmt.Printf("关键词位命中: %s\n", seq)
				return make([]DocStore, 0), nil
			} else {
				return nil, err
			}
		} else {
			vsc := len(v) / 16
			vids := mapset.NewSetWithSize[uuid.UUID](vsc)
			for i := 0; i < vsc; i++ {
				if id, err := uuid.FromBytes(v[i*16 : (i+1)*16]); err != nil {
					return result, err
				} else {
					vids.Add(id)
				}
			}
			fmt.Printf("命中关键词：%s  关联数: %d\n", seq, vsc)

			if ids.IsEmpty() {
				ids = vids
			} else {
				ids = ids.Intersect(vids)
			}
		}
	}

	// 获取文档内容
	it := ids.Iterator()
	var doc DocContent
	for id := range it.C {
		d, err := DocDb.Get(id[:], &opt.ReadOptions{})
		if err != nil {
			return result, err
		}
		if err := json.Unmarshal(d, &doc); err != nil {
			return result, err
		}
		result = append(result, DocStore{
			DocContent: doc,
			ID:         id,
		})
	}

	return result, nil
}
