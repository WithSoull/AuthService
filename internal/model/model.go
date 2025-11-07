package model

type UserInfo struct {
	UserId int64
	Email  string
}

func (ui UserInfo) GetUserID() int64 {
	return ui.UserId
}

func (ui UserInfo) GetEmail() string {
	return ui.Email
}
