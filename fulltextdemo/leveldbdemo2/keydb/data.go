package keydb

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataTask struct {
}

var DataTaskQueue chan DataTask

var Mysql *gorm.DB

func InitMysql() error {
	dns := os.Getenv("DB_DNS")

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return err
	}
	Mysql = db

	go func() {
		DataTaskQueue = make(chan DataTask)
		fmt.Println("start mysql task queue.")
		for {
			<-DataTaskQueue
			if err := ImportData(); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("final mysql task.")
		}
	}()

	return nil
}

func initStartId() any {
	keyType := os.Getenv("DB_TASK_KEY_TYPE")
	switch keyType {
	case "uint32":
		return uint32(0)
	default:
		fmt.Printf("mysql task: unknown startId %s\n", keyType)
		return int(0)
	}
}

func pickEndId(endRow map[string]any, key string) any {
	endKey := endRow[key]
	fmt.Printf("mysql task end key: %v\n", endKey)
	switch endKey.(type) {
	case string:
		fmt.Printf("endKey: is string\n")
	case uint32:
		fmt.Printf("endKey: is uint32\n")
	default:
		fmt.Printf("endKey: type %v\n", reflect.TypeOf(endKey))
	}
	return endKey
}

func ImportData() error {
	table := os.Getenv("DB_TASK_TABLE")
	key := os.Getenv("DB_TASK_KEY")
	limit := os.Getenv("DB_TASK_LIMIT")
	plainKey := os.Getenv("DB_TASK_PLAIN_KEY")
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s > ? ORDER BY %s LIMIT %s", table, key, key, limit)
	fmt.Printf("mysql task: %s\n", sql)

	startId := initStartId()
	for {
		rows := []map[string]any{}
		if err := Mysql.Raw(sql, startId).Scan(&rows).Error; err != nil {
			return err
		}
		count := len(rows)
		if count <= 0 {
			break
		}
		endRow := rows[count-1]
		endId := pickEndId(endRow, key)
		fmt.Printf("mysql task fetch row count: %d endId: %v\n", count, endId)
		startId = endId

		docs := make([]DocContent, count)
		for i, row := range rows {
			if plain, ok := row[plainKey].(string); ok {
				docs[i] = DocContent{
					Plain:   plain,
					Content: row,
				}
			} else {
				fmt.Printf("mysql task invalid plainKey(%s) for %v\n", plainKey, row)
			}
		}
		if _, err := AddBatch(docs); err != nil {
			fmt.Printf("mysql task add batch error: %v\n", err)
		}
	}

	return nil
}
