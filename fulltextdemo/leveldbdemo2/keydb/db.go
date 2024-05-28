package keydb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

var DocDb *leveldb.DB
var IndexDb *leveldb.DB

func InitDb() func() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cur: %s\n", pwd)
	docDbDir := filepath.Join(pwd, "doc.ldb")
	indexDbDir := filepath.Join(pwd, "index.ldb")

	fmt.Printf("docDb dir: %s\n", docDbDir)
	fmt.Printf("indexDb dir: %s\n", indexDbDir)

	// doc
	docDb, err := leveldb.OpenFile(docDbDir, nil)
	if err != nil {
		log.Fatal(err)
	}
	DocDb = docDb

	// index
	indexDb, err := leveldb.OpenFile(indexDbDir, nil)
	if err != nil {
		log.Fatal(err)
	}
	IndexDb = indexDb

	return func() {
		indexDb.Close()
		docDb.Close()
	}
}
