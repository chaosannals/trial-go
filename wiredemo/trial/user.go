package trial

import (
	"fmt"
	"time"
)

type User struct {
	Id          int64
	Account     string
	Nickname    string
	LastLoginAt time.Time
}

func NewUser(account string) *User {
	return &User{
		Id: 1,
		Account: account,
		Nickname: fmt.Sprintf("%s - %d", account, 1),
		LastLoginAt: time.Now(),
	}
}