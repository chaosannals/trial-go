package tests

import (
	"encoding/json"

	"github.com/revel/revel/testing"
)

type VisitTest struct {
	testing.TestSuite
}

func (t *VisitTest) TestList() {
	t.Get("/visit/list?p=2&s=30")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	var data map[string]interface{}
	json.Unmarshal(t.ResponseBody, &data)
	t.AssertEqual(data["p"], "2")
	t.AssertEqual(data["s"], "30")
	_, ok := data["rows"]
	t.Assert(ok)
}

func (t *VisitTest) TestInfo() {
	t.Get("/visit/123")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	var data map[string]interface{}
	json.Unmarshal(t.ResponseBody, &data)
	t.AssertEqual(data["id"], "123")
}
