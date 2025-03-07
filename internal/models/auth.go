package models

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type VerifyOTPReq struct {
	Token string `json:"token" form:"token"`
	Otp   string `json:"otp" form:"otp"`
}

type VerifyOTPRes struct {
	TTL int `json:"ttl"`
}

type ResendOTPReq struct {
	Token string `json:"token" form:"token"`
}
