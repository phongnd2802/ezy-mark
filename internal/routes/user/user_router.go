package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/controllers"
	"github.com/phongnd2802/ezy-mark/internal/middlewares"
)

type UserRouter struct{}

func (r *UserRouter) InitUserRouter(Router fiber.Router) {
	userRouter := Router.Group("/user")
	userRouter.Use(middlewares.AuthenticationMiddleware())
	{
		userRouter.Patch("/update-info", controllers.UserInfo.UpdateUserProfile)
		userRouter.Get("/get-info", controllers.UserInfo.GetUserProfile)
	}
}