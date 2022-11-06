package trial

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormDb() (db *gorm.DB, err error) {
	cs := "root:password@tcp(localhost:3306)/exert?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(cs), &gorm.Config{})
}
