package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/config"
	"github.com/phongnd2802/daily-social/internal/controllers/account"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/global"
	"github.com/phongnd2802/daily-social/internal/response"
	"github.com/phongnd2802/daily-social/internal/services"
	"github.com/phongnd2802/daily-social/internal/services/impl"
	"github.com/phongnd2802/daily-social/internal/worker"
	"github.com/phongnd2802/daily-social/pkg/email"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

	app.Get("/", func(c *fiber.Ctx) error {
		return response.SuccessResponse(c, response.ErrCodeSuccess, "Success")
	})

	app.Post("/api/v1/auth/signup", account.Auth.Register)
	app.Post("/api/v1/auth/verify-otp", account.Auth.VerifyOTP)

	app.Get("/api/v1/auth/verify-otp", account.Auth.GetTTLOtp)

	app.Listen("127.0.0.1:8000")
}


func runTaskProcessor(reditOpt asynq.RedisClientOpt, sender email.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(reditOpt, sender)
	log.Info().Msg("Start Task Processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Msg("failed to start task processor")
	}
}
