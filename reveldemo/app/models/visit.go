package models

import "time"

type VisitModel struct {
	Id       uint64    `db:"id, primarykey, autoincrement"`
	Name     string    `db:"name, size:100"`
	Pass     string    `db:"pass, size:64"`
	CreateAt time.Time `db:"create_at"`
}
