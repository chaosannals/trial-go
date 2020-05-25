package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //sqlite
)

//Store 存储
type Store struct {
	gorm.Model
	ID      string `gorm:"primary_key;AUTO_INCREMENT;comment:'APP ID'"`
	Title   string `gorm:"not null;size:255;comment:'标题';"`
	Content string `gorm:"not null;type:LONGTEXT;comment:'标题';"`
	TypeID  uint64 `gorm:"comment:'类型ID';"`
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
