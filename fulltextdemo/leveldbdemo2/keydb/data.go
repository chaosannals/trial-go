package keydb

import (
	"fmt"
	"log"
	"os"

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

func ImportData() error {
	table := os.Getenv("DB_TASK_TABLE")
	key := os.Getenv("DB_TASK_KEY")
	limit := os.Getenv("DB_TASK_LIMIT")
	// contentKey := os.Getenv("DB_TASK_CONTENT_KEY")
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s > ? ORDER BY %s LIMIT %s", table, key, key, limit)
	result := []map[string]any{}
	startId := 0
	fmt.Printf("mysql task: %s", sql)
	for {
		if err := Mysql.Raw(sql, startId).Scan(&result).Error; err != nil {
			return err
		}
		count := len(result)
		if count <= 0 {
			break
		}
		endItem := result[count-1]
		if endId, ok := endItem[key].(int); ok {
			startId = endId
		}
		fmt.Printf("mysql task fetch row count: %d\n", count)

	}

	return nil
}
