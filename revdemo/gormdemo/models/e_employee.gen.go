// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameEEmployee = "e_employee"

// EEmployee mapped from table <e_employee>
type EEmployee struct {
	ID          uint64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id"`
	Account     string     `gorm:"column:account;type:varchar(30);not null;uniqueIndex:ACCOUNT_UNIQUE,priority:1" json:"account"` // 账号
	Password    *[]byte    `gorm:"column:password;type:binary(32)" json:"password"`
	Nickname    *string    `gorm:"column:nickname;type:varchar(30)" json:"nickname"`                                     // 昵称
	CreatedAt   time.Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	RemovedAt   *time.Time `gorm:"column:removed_at;type:datetime" json:"removed_at"`                                    // 删除时间
	LastLoginAt *time.Time `gorm:"column:last_login_at;type:datetime" json:"last_login_at"`                              // 最后登录时间
}

// TableName EEmployee's table name
func (*EEmployee) TableName() string {
	return TableNameEEmployee
}
