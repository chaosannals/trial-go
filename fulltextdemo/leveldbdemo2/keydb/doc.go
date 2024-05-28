package keydb

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type DocContent struct {
	Content   string  `json:"content" form:"content"`
	IntInfo   int32   `json:"intInfo" form:"intInfo"`
	FloatInfo float64 `json:"floatInfo" form:"floatInfo"`
}

func (doc *DocContent) Encode() ([]byte, error) {
	return json.Marshal(doc)
}

func (doc *DocContent) InsertAndCut() (uuid.UUID, []string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return id, nil, err
	}

	de, err := doc.Encode()
	if err != nil {
		return id, nil, err
	}

	// seqs := gse.ToString(keydb.Seg.Segment([]byte(doc.Content)), true)
	seqs := Seg.CutSearch(doc.Content, true)

	for _, s := range seqs {
		k := []byte(s)
		r := id[:]
		if v, err := DocDb.Get(k, &opt.ReadOptions{}); err != nil {
			if err == leveldb.ErrNotFound {
				fmt.Printf("插入 %v\n", r)
			} else {
				fmt.Printf("插入错误 %v\n", err)
			}
		} else {
			r = append(r, v...)
		}
		DocDb.Put([]byte(s), r, &opt.WriteOptions{})
	}

	err = DocDb.Put(id[:], de, &opt.WriteOptions{})

	return id, seqs, err
}

type DocIndex = []uuid.UUID
