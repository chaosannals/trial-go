package logic

import (
	"context"

	"sqlitedemo/internal/svc"
	"sqlitedemo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SqlitedemoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSqlitedemoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SqlitedemoLogic {
	return &SqlitedemoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SqlitedemoLogic) Sqlitedemo(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
