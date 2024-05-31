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
	Plain   string         `json:"plain" form:"plain"`
	Content map[string]any `json:"content" form:"content"`
}
type DocStore struct {
	DocContent
	ID uuid.UUID `json:"id" form:"id"`
}

type DocIndex = []uuid.UUID

func (doc *DocContent) CutIndex(id uuid.UUID) ([]string, error) {
	// 索引
	// seqs := gse.ToString(keydb.Seg.Segment([]byte(doc.Content)), true)
	seqs := Seg.CutSearch(doc.Plain, true)
	batch := new(leveldb.Batch)
	for _, s := range seqs {
		k := []byte(s)
		r := id[:]
		if v, err := IndexDb.Get(k, &opt.ReadOptions{}); err != nil {
			if err == leveldb.ErrNotFound {
				fmt.Printf("插入 %v\n", r)
			} else {
				fmt.Printf("插入错误 %v\n", err)
				return seqs, err
			}
		} else {
			r = append(r, v...)
		}
		batch.Put(k, r)
	}
	err := IndexDb.Write(batch, &opt.WriteOptions{})
	return seqs, err
}

type AddResult struct {
	ID   uuid.UUID `json:"id" form:"id"`
	Seqs []string  `json:"seqs" form:"seqs"`
}

func (doc *DocContent) InsertAndCut() (*AddResult, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	de, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// 索引
	seqs, err := doc.CutIndex(id)
	if err != nil {
		return nil, err
	}

	// 文档
	err = DocDb.Put(id[:], de, &opt.WriteOptions{})

	return &AddResult{ID: id, Seqs: seqs}, err
}

func AddBatch(docs []DocContent) ([]AddResult, error) {
	result := make([]AddResult, len(docs))
	batch := new(leveldb.Batch)
	for i, doc := range docs {
		id, err := uuid.NewUUID()
		if err != nil {
			return result, err
		}
		seqs, err := doc.CutIndex(id)
		if err != nil {
			return result, err
		}
		de, err := json.Marshal(doc)
		if err != nil {
			return result, err
		}
		batch.Put(id[:], de)
		result[i] = AddResult{
			ID:   id,
			Seqs: seqs,
		}
	}

	err := DocDb.Write(batch, &opt.WriteOptions{})
	return result, err
}

func ParseUUIDSet(v []byte) (mapset.Set[uuid.UUID], error) {
	vsc := len(v) / 16
	vids := mapset.NewSetWithSize[uuid.UUID](vsc)
	for i := 0; i < vsc; i++ {
		if id, err := uuid.FromBytes(v[i*16 : (i+1)*16]); err != nil {
			return vids, err
		} else {
			vids.Add(id)
		}
	}
	return vids, nil
}

type QueryParam struct {
	Plain   string
	GroupBy *string
}

// 切词 AND 逻辑
func Query(param *QueryParam) ([]DocStore, error) {
	result := make([]DocStore, 0)

	fmt.Printf("检索： %s\n", param.Plain)
	seqs := Seg.CutSearch(param.Plain, true)

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
			vids, err := ParseUUIDSet(v)
			if err != nil {
				return result, err
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
