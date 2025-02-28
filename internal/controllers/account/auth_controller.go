package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/daily-social/internal/dtos"
	"github.com/phongnd2802/daily-social/internal/response"
	"github.com/phongnd2802/daily-social/internal/services"
)

type authController struct{}

var Auth = new(authController)


func (c *authController) Login(ctx *fiber.Ctx) error {
	params := new(dtos.LoginRequest)
	if err := ctx.BodyParser(params); err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}

	code, data, err := services.AuthService().Login(ctx.UserContext(), params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, code, data)
}

func (c *authController) Register(ctx *fiber.Ctx) error {
	params := new(dtos.RegisterRequest)
	if err := ctx.BodyParser(params); err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}
	/// Call Auth Service
	code, data, err := services.AuthService().Register(ctx.UserContext(), params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}	
	return response.SuccessResponse(ctx, code, data)
}

func (c *authController) VerifyOTP(ctx *fiber.Ctx) error {
	params := new(dtos.VerifyOTPReq)
	if err := ctx.BodyParser(params); err != nil {
		return response.ErrorResponse(ctx, response.ErrCodeInvalidParams, err)
	}
	code, data, err := services.AuthService().VerifyOTP(ctx.UserContext(), params)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, code, data)
}

func (c *authController) GetTTLOtp(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return response.ErrorResponse(ctx, response.ErrCodeExpiredSession, nil)
	}
	code, data, err := services.AuthService().GetTTLOtp(ctx.UserContext(), token)
	if err != nil {
		return response.ErrorResponse(ctx, code, err)
	}
	return response.SuccessResponse(ctx, code, data)
}