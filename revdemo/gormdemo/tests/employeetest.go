package tests

import (
	"github.com/revel/revel/testing"
)

type EmployeeTest struct {
	testing.TestSuite
}

func (t *EmployeeTest) TestFind() {
	t.Get("/employee/find")
	t.AssertOk()
	t.AssertContains("rows")
}