package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/handlers"
	"github.com/redis/go-redis/v9"
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
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loding .env file")
	}

	connPool, err := pgxpool.New(context.Background(), os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatalf("unable to connect database %v", err)
	}

	err = connPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("failed ping database %v", err)
	}

	log.Println("Connected Database Success")

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	err = client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	log.Println("Connected Redis Success")

	store := db.NewStore(connPool)
	cache := cache.NewRedisClient(client)

	handler := handlers.NewHandler(store, cache)
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))
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
