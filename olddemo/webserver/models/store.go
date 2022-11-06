package models

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
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
	db, e := gorm.Open(sqlite.Open("store.db"), &gorm.Config{})
	if e != nil {
		panic(e.Error())
	}
	// 老版初始化
	// db.CreateTable(&Store{})
	
	db.AutoMigrate(&Store{})
	return db
}
