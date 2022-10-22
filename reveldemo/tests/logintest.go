package tests

import (
	"github.com/revel/revel/testing"
)

type LoginTest struct {
	testing.TestSuite
}

func (t *LoginTest) TestIndex() {
	t.Get("/login.html")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}
