package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //sqlite
)

//Store 存储
type Store struct {
	gorm.Model
	ID      string // ID
	Title   string // 标题
	Content string // 内容
	TypeID  uint64 // 类型ID
}

//InitStore 初始化
func InitStore() *gorm.DB {
	db, e := gorm.Open("sqlite3", "store.db")
	if e != nil {
		panic(e.Error())
	}
	db.CreateTable(&Store{})
	return db
}
