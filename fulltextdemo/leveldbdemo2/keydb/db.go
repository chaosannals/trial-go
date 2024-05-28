package keydb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

var LDB *leveldb.DB

func InitDb() func() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cur: %s\n", pwd)
	dbDir := filepath.Join(pwd, "demo.ldb")
	fmt.Printf("db dir: %s\n", dbDir)
	db, err := leveldb.OpenFile(dbDir, nil)
	if err != nil {
		log.Fatal(err)
	}
	LDB = db
	return func() {
		db.Close()
	}
}
