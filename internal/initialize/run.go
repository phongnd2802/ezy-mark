package initialize

import "github.com/gofiber/fiber/v2"

func Run() *fiber.App {
	loadConfig()
	initDatabase()
	initRedis()
	initServiceInterfaces()

	app := initRouter()

	return app
}