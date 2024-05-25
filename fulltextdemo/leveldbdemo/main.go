package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-faker/faker/v4"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
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
	defer db.Close()

	batch := new(leveldb.Batch)
	for i := 0; i < 1000; i++ {
		uuid := faker.UUIDDigit()
		name := faker.ChineseName()
		batch.Put([]byte(uuid), []byte(name))
	}
	if err := db.Write(batch, nil); err != nil {
		log.Fatal(err)
	}

	k1 := []byte("aaa")
	if err := db.Put(k1, []byte("123"), nil); err != nil {
		log.Fatal(err)
	}
	if v1, err := db.Get(k1, nil); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("v1: %v\n", string(v1))
	}

	if err := db.Delete(k1, nil); err != nil {
		log.Fatal(err)
	}
}
