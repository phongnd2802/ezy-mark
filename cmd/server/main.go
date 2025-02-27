package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/config"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/global"
	"github.com/phongnd2802/daily-social/internal/handlers"
	"github.com/phongnd2802/daily-social/internal/worker"
	"github.com/phongnd2802/daily-social/pkg/email"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}

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

	go runTaskProcessor(redisOpt, sender)

	handler := handlers.NewHandler(store, cache, distributor)
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SessionSecret))))
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())
	e.Static("/public", "public")
	e.GET("/foo", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"code":    200,
			"message": "success",
			"result":  "FOO",
		})
	})

	e.GET("/", handler.HandleHome)
	e.POST("/", handler.HandleHome)

	e.GET("/signin", handler.HandleLogin)
	e.POST("/signin", handler.HandleLogin)

	e.GET("/signup", handler.HandleRegister)
	e.POST("/signup", handler.HandleRegister)

	e.GET("/verify-otp", handler.HandleVerifyOTP)
	e.POST("/verify-otp", handler.HandleVerifyOTP)

	e.GET("/forgot-password", handler.HandleForgotPassword)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := e.Start("127.0.0.1:1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Server gracefully stopped")
}

func runTaskProcessor(reditOpt asynq.RedisClientOpt, sender email.EmailSender) {
	taskProcessor := worker.NewRedisTaskProcessor(reditOpt, sender)
	log.Info().Msg("Start Task Processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Msg("failed to start task processor")
	}
}
