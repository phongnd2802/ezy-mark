package models

import "mime/multipart"

type UserProfile struct {
	UserNickname string  `json:"user_nickname" form:"user_nickname"`
	UserFullname string  `json:"user_fullname" form:"user_fullname"`
	UserMobile   string  `json:"user_mobile" form:"user_mobile"`
	UserGender   *string `json:"user_gender" form:"user_gender"`
	UserBirthday string  `json:"user_birthday" form:"user_birthday"`
}

type UpdateProfileUserReq struct {
	UserId int64
	UserProfile
	UserAvatar *multipart.FileHeader
	SubToken   string
}

type UpdateProfileUserRes struct {
	UserAvatar string `json:"user_avatar"`
	UserProfile
}

type ChangePassword struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type GetProfileParams struct {
	UserId   int64
	SubToken string
}
