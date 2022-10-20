package controller

import (
	"context"
	v1 "gfdemo/api/v1"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	Login = cLogin{}
)

type cLogin struct {
}

func (c *cLogin) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln("{ \"code\": 1 }")
	return
}
