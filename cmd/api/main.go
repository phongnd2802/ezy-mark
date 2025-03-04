package main

import (
	"context"
	"github.com/phongnd2802/daily-social/internal/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/config"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/global"
	"github.com/phongnd2802/daily-social/internal/middlewares"
	"github.com/phongnd2802/daily-social/internal/pkg/email"
	"github.com/phongnd2802/daily-social/internal/response"
	"github.com/phongnd2802/daily-social/internal/services"
	"github.com/phongnd2802/daily-social/internal/services/impl"
	"github.com/phongnd2802/daily-social/internal/worker"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/swagger"
	_ "github.com/phongnd2802/daily-social/docs"
)

// @title Daily Social Fiber API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name Philip Nguyen
// @contact.email duyphong0280@gmail.com
// @license.name MIT
// @license.url https://github.com/phongnd2802/daily-social/blob/main/LICENSE
// @host localhost:8000
// @BasePath /api/v1
func main() {
	config, err := config.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	if config.Mode == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msgf("unable to connect database %v", err)
	}
	defer connPool.Close()

	err = connPool.Ping(context.Background())
	if err != nil {
		log.Fatal().Msgf("failed ping database %v", err)
	}

	log.Info().Msg("Connected Database Success")

	global.ConnPool = connPool

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	defer client.Close()

	err = client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Info().Msg("Connected Redis Success")

	global.Rdb = client

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddr,
	}

	store := db.NewStore()
	cache := cache.NewRedisClient()
	distributor := worker.NewRedisTaskDistributor(redisOpt)
	sender := email.NewGmailSender(
		config.SenderName,
		config.SenderEmail,
		config.SenderPassword,
	)

	services.InitAuthService(impl.NewAuthServiceImpl(store, cache, distributor))

	go runTaskProcessor(redisOpt, sender)
	app := fiber.New()

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	app.Use(compress.New())
	app.Use(middlewares.RequestMetadataMiddleware())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/api/v1/foo", func(c *fiber.Ctx) error {
		return response.SuccessResponse(c, response.ErrCodeSuccess, "Success")
	})

	app.Post("/api/v1/auth/signup", controllers.Auth.Register)
	app.Post("/api/v1/auth/verify-otp", controllers.Auth.VerifyOTP)
	app.Post("/api/v1/auth/login", controllers.Auth.Login)

	app.Get("/api/v1/auth/verify-otp", controllers.Auth.GetTTLOtp)
	app.Post("/api/v1/auth/resend-otp", controllers.Auth.ResendOTP)
	app.Listen("127.0.0.1:8000")
}

func runTaskProcessor(reditOpt asynq.RedisClientOpt, sender email.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(reditOpt, sender)
	log.Info().Msg("Start Task Processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Msg("failed to start task processor")
	}
}
