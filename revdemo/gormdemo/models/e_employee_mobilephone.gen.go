// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameEEmployeeMobilephone = "e_employee_mobilephone"

// EEmployeeMobilephone mapped from table <e_employee_mobilephone>
type EEmployeeMobilephone struct {
	EmployeeID  uint64    `gorm:"column:employee_id;type:bigint unsigned;primaryKey" json:"employee_id"` // 员工ID
	Mobilephone string    `gorm:"column:mobilephone;type:varchar(15);primaryKey" json:"mobilephone"`     // 手机号码
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName EEmployeeMobilephone's table name
func (*EEmployeeMobilephone) TableName() string {
	return TableNameEEmployeeMobilephone
}
