package entities

import "time"

// UserInfo .
type UserInfo struct {
	UID        int        `table:"userinfo" column:"uid"`
	UserName   string     `column:"username"`
	DepartName string     `column:"departname"`
	CreateAt   *time.Time `column:"created"`
}

// NewUserInfo .
func NewUserInfo(u UserInfo) *UserInfo {
	if len(u.UserName) == 0 {
		panic("UserName shold not null!")
	}
	if u.CreateAt == nil {
		t := time.Now()
		u.CreateAt = &t
	}
	return &u
}
