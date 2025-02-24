package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phongnd2802/daily-social/internal/handlers"
)

func main() {
	e := echo.New()
	// e.Use(middleware.WithCSP(middleware.CSPConfig{
	// 	ScriptSrc: []string{"/public/js/htmx.min.js"},
	// }))
	e.Static("/public", "public")
	e.GET("/foo", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"code":    200,
			"message": "success",
			"result":  "FOO",
		})
	})

	e.GET("/", handlers.HandleHome)

	e.GET("/signin", handlers.HandleLogin)
	e.POST("/signin", handlers.HandleLogin)

	e.GET("/signup", handlers.HandleRegister)
	e.POST("/signup", handlers.HandleRegister)

	e.GET("/verify-otp", handlers.HandleVerifyOTP)

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
