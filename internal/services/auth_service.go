package services

import (
	"context"

	"github.com/phongnd2802/daily-social/internal/dtos"
)

type iAuthService interface {
	Register(ctx context.Context, params *dtos.RegisterRequest) (code int, res *dtos.RegisterResponse, err error)
	Login(ctx context.Context, params *dtos.LoginRequest) (code int, res *dtos.LoginResponse, err error)
	VerifyOTP(ctx context.Context, params *dtos.VerifyOTPReq) (code int, res *dtos.VerifyOTPRes, err error)

	ResendOTP(ctx context.Context, params *dtos.ResendOTPReq) (code int, err error)
	GetTTLOtp(ctx context.Context, token string) (code int, res *dtos.VerifyOTPRes, err error)
}

var localAuthService iAuthService

func AuthService() iAuthService {
	if localAuthService == nil {
		panic("iAuthService interface not implemented yet")
	}
	return localAuthService
}

func InitAuthService(i iAuthService) {
	localAuthService = i
}
