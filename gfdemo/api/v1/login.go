package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type LoginReq struct {
	g.Meta `path:"/login" method:"post" summary:"login api"`
}

type LoginRes struct {
	g.Meta `mime:"application/json" example:"{ \"code\": 0 }"`
}
