package models

import "time"

type UserProfile struct {
	UserNickname string    `json:"user_nickname" form:"user_nickname"`
	UserFullname string    `json:"user_fullname" form:"user_fullname"`
	UserAvatar   string    `json:"user_avatar" form:"user_avatar"`
	UserMobile   string    `json:"user_mobile" form:"user_mobile"`
	UserGender   bool      `json:"user_gender" form:"user_gender"`
	UserBirthday time.Time `json:"user_birthday" form:"user_birthday"`
}

type UpdateProfileUserReq struct {
	UserProfile
}

type UpdateProfileUserRes struct {
	UserProfile
}

type ChangePassword struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

