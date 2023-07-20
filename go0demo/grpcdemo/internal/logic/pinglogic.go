package logic

import (
	"context"

	"grpcdemo/grpcdemo"
	"grpcdemo/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *grpcdemo.Request) (*grpcdemo.Response, error) {
	// 此文件代码虽然是生成文件，但是通过 server 隔离开的，所以 GRPC 的生成物没有和此文件混合。

	return &grpcdemo.Response{
		Pong:"pong",
	}, nil
}
