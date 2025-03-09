package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/controllers"
	"github.com/phongnd2802/ezy-mark/internal/middlewares"
)

type ShopRouter struct{}

func (r *ShopRouter) InitShopRouter(Router fiber.Router) {
	shopRouter := Router.Group("/shop")
	shopRouter.Use(middlewares.AuthenticationMiddleware())
	shopRouter.Use(middlewares.RBAC([]string{"customer", "shop"}))
	{
		shopRouter.Post("/register", controllers.ShopUser.RegisterShop)
	}
}
