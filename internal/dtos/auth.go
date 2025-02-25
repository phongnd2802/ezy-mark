package dtos

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
type RegisterRequest struct {
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

type VerifyOTPReq struct {
	Token string `json:"token" form:"token"`
	Otp   string `json:"otp" form:"otp"`
}
