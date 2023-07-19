package handler

import (
	"net/http"

	"apidemo/internal/logic"
	"apidemo/internal/svc"
	"apidemo/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ApidemoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewApidemoLogic(r.Context(), svcCtx)
		resp, err := l.Apidemo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
