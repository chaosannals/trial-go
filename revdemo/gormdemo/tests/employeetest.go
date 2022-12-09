package tests

import (
	"bytes"
	"encoding/json"
	"gormdemo/models"

	"github.com/go-faker/faker/v4"
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

func (t *EmployeeTest) TestInfo() {
	t.Get("/employee/info/1")
	t.AssertOk()
	t.AssertContains("info")
}

func (t *EmployeeTest) TestAdd() {
	var m models.EEmployee
	err := faker.FakeData(&m)
	t.AssertEqual(err, nil)
	b, err := json.Marshal(&m)
	t.AssertEqual(err, nil)
	t.Put("/employee/add", "application/json", bytes.NewBuffer(b))
	t.AssertOk()
}

func (t *EmployeeTest) TestEdit() {
	m := &models.EEmployee{
		ID:       2,
		Account:  "12312",
		Nickname: nil,
	}
	b, err := json.Marshal(m)
	t.AssertEqual(err, nil)
	t.Post("/employee/edit", "application/json", bytes.NewReader(b))
	t.AssertOk()
}
