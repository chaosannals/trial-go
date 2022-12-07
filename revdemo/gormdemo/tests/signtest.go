package tests

import (
	"strings"

	"github.com/revel/revel/testing"
)

type SignTest struct {
	testing.TestSuite
}

func (t *SignTest) TestLogin() {
	t.Post("/login", "application/json", strings.NewReader("{}"))
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}
