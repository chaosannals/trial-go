package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sqlitedemo/internal/logic"
	"sqlitedemo/internal/svc"
	"sqlitedemo/internal/types"
)

func SqlitedemoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSqlitedemoLogic(r.Context(), svcCtx)
		resp, err := l.Sqlitedemo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
