package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/controllers"
	"github.com/phongnd2802/ezy-mark/internal/middlewares"
)

type AuthRouter struct{}

func (r *AuthRouter) InitAuthRoute(Router fiber.Router) {

	authPublic := Router.Group("/auth")
	{
		authPublic.Post("/signup", controllers.Auth.Register)
		authPublic.Post("/signin", controllers.Auth.Login)

		authPublic.Post("/resend-otp", controllers.Auth.ResendOTP)
		authPublic.Post("/verify-otp", controllers.Auth.VerifyOTP)
		authPublic.Get("/verify-otp", controllers.Auth.GetTTLOtp)
	}

	authPrivate := Router.Group("/auth")
	authPrivate.Use(middlewares.AuthenticationMiddleware())
	{
		authPrivate.Post("/logout", nil)
	}
}
