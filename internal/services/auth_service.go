package services

import (
	"context"

	"github.com/phongnd2802/ezy-mark/internal/models"
)

type IAuthService interface {
	// Login, Register, Logout
	Register(ctx context.Context, params *models.RegisterRequest) (code int, res *models.RegisterResponse, err error)
	Login(ctx context.Context, params *models.LoginRequest) (code int, res *models.LoginResponse, err error)
	Logout(ctx context.Context, sessionId string) (code int, err error)

	// Verify OTP
	VerifyOTP(ctx context.Context, params *models.VerifyOTPReq) (code int, res *models.VerifyOTPRes, err error)
	ResendOTP(ctx context.Context, params *models.ResendOTPReq) (code int, err error)
	GetTTLOtp(ctx context.Context, token string) (code int, res *models.VerifyOTPRes, err error)

	// Handle Refresh Token
	RefreshToken(ctx context.Context, params *models.RefreshTokenParams) (code int, res *models.LoginResponse, err error)

	// Handle Forgot Password
	// TODO...
}

var localAuthService IAuthService

func AuthService() IAuthService {
	if localAuthService == nil {
		panic("IAuthService interface not implemented yet")
	}
	return localAuthService
}

func InitAuthService(i IAuthService) {
	localAuthService = i
}
