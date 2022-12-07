// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entities

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:                   db,
		EEmployee:            newEEmployee(db, opts...),
		EEmployeeMobilephone: newEEmployeeMobilephone(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	EEmployee            eEmployee
	EEmployeeMobilephone eEmployeeMobilephone
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:                   db,
		EEmployee:            q.EEmployee.clone(db),
		EEmployeeMobilephone: q.EEmployeeMobilephone.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.clone(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.clone(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:                   db,
		EEmployee:            q.EEmployee.replaceDB(db),
		EEmployeeMobilephone: q.EEmployeeMobilephone.replaceDB(db),
	}
}

type queryCtx struct {
	EEmployee            *eEmployeeDo
	EEmployeeMobilephone *eEmployeeMobilephoneDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		EEmployee:            q.EEmployee.WithContext(ctx),
		EEmployeeMobilephone: q.EEmployeeMobilephone.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	return &QueryTx{q.clone(q.db.Begin(opts...))}
}

type QueryTx struct{ *Query }

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
