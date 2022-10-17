package basis

import (
	"time"
)

// 时间格式结构
type MyTime struct {
	time.Time
}

// 格式化时间。
func (i MyTime) MarshalJSON() ([]byte, error) {
	text := i.Format("\"2006-01-02 15:04:05\"")
	return []byte(text), nil
}

// 解析时间。
func (i *MyTime) UnmarshalJSON(data []byte) error {
	text := string(data)
	t, err := time.Parse("\"2006-01-02 15:04:05\"", text)
	*i = MyTime{t}
	return err
}

// 访问令牌
type AccessToken struct {
	CompanyId uint32 `json:"companyId"`
	CreateAt  MyTime `json:"createAt"`
}

// 加密
func (i *AccessToken) Encrypt(key []byte) (string, error) {
	return Encrypt(key, i)
}

// 从字符串生成访问令牌。
func FromString(key []byte, text string) (*AccessToken, error) {
	r := &AccessToken{}
	err := Decrypt(key, text, r)
	return r, err
}
