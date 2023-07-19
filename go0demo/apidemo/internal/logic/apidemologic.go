package logic

import (
	"context"

	"apidemo/internal/svc"
	"apidemo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApidemoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApidemoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApidemoLogic {
	return &ApidemoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApidemoLogic) Apidemo(req *types.Request) (resp *types.Response, err error) {
	// 此处添加代码，本文件是生成代码和手写代码的混合文件。
	resp = new(types.Response)
	resp.Message = req.Name
	return
}
