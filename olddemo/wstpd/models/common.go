package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 公司软件信息
type CompanySoftware struct {
	CompanyId       uint32    `gorm:"column:company_id;primaryKey;comment:'公司ID'"`
	Appkey          string    `gorm:"column:appkey;comment:'APPKEY'"`
	Appsecret       string    `gorm:"column:appsecret;comment:'APPSECRET'"`
	SoftwareType    uint8     `gorm:"column:software_type;comment:'软件类型'"`
	SoftwareVersion string    `gorm:"column:software_version;comment:'软件版本'"`
	Contact         string    `gorm:"column:contact;comment:'联系人'"`
	Telephone       string    `gorm:"column:telephone;comment:'电话'"`
	Mobilephone     string    `gorm:"column:mobilephone;comment:'手机'"`
	SyncCount       int32     `gorm:"column:sync_count;comment:'数量'"`
	SyncTime        time.Time `gorm:"column:sync_time;comment:'时间'"`
	// LocalIp string
	// UseIp string
	CreateTime time.Time `gorm:"column:create_time;comment:'创建时间'"`
	Remark     string    `gorm:"column:remark;comment:'注释'"`
}

// 表名
func (i *CompanySoftware) TableName() string {
	return "tg_company_software"
}


// 初始化公共数据库的连接
func OpenCommonDb() (db *gorm.DB, err error) {
	name := os.Getenv("DB_NAME")
	if len(name) == 0 {
		return nil, errors.New("main db empty DB_NAME")
	}
	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		return nil, errors.New("main db empty DB_HOST")
	}
	port := os.Getenv("DB_PORT")
	if len(port) == 0 {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if len(user) == 0 {
		return nil, errors.New("main db empty DB_USER")
	}
	pass := os.Getenv("DB_PASS")
	if len(pass) == 0 {
		return nil, errors.New("main db empty DB_PASS")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
