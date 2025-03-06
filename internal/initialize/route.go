package initialize

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/phongnd2802/ezy-mark/internal/middlewares"
	"github.com/phongnd2802/ezy-mark/internal/routes"
)

func initRouter() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(recover.New())
	v1.Use(requestid.New())
	v1.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	v1.Use(compress.New())
	v1.Use(middlewares.RequestMetadataMiddleware())

	userRouter := routes.RouterGroupApp.User
	userRouter.InitAuthRoute(v1)

	return app
}
