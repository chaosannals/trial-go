package trial

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormDb(conf *Conf) (db *gorm.DB, err error) {
	cs := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DbUser,
		conf.DbPass,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
	)
	return gorm.Open(mysql.Open(cs), &gorm.Config{})
}
