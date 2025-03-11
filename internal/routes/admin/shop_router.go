package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phongnd2802/ezy-mark/internal/controllers"
	"github.com/phongnd2802/ezy-mark/internal/middlewares"
)

type ShopRouter struct{}

func (r *ShopRouter) InitShopRouter(Router fiber.Router) {
	shopRouter := Router.Group("/admin/shop")
	shopRouter.Use(middlewares.AuthenticationMiddleware())
	shopRouter.Use(middlewares.RBAC([]string{"admin"}))
	{
		shopRouter.Patch("/approve", controllers.ShopAdmin.ApproveShop)
	}
}
