package util

import (
	"encoding/json"
	"net/http"

	"github.com/revel/revel"
)

type MyResult struct {
	Content interface{}
}

func (r MyResult) Apply(req *revel.Request, resp *revel.Response) {
	var b []byte
	var err error

	if revel.Config.BoolDefault("results.pretty", false) {
		b, err = json.MarshalIndent(r.Content, "", "  ")
	} else {
		b, err = json.Marshal(r.Content)
	}

	if err != nil {
		revel.ErrorResult{Error: err}.Apply(req, resp)
		return
	}

	resp.WriteHeader(http.StatusOK, "application/javascript; charset=utf-8")
	if _, err = resp.GetWriter().Write(b); err != nil {
		revel.AppLog.Error("Apply: Response write failed", "error", err)
	}
}