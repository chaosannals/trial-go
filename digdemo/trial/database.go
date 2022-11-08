package trial

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/chaosannals/trial-go-digdemo/entities"
	"github.com/chaosannals/trial-go-digdemo/models"
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
	if db, err = gorm.Open(mysql.Open(cs), &gorm.Config{}); err != nil {
		return
	}

	// 启用生成的模型和查询代码
	entities.Use(db)

	// CodeFirst
	err = db.AutoMigrate(
		models.EEmployee{},
		models.EEmployeeMobilephone{},
	)
	return
}
