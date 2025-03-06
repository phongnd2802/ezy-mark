package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/phongnd2802/ezy-mark/internal/initialize"

	"github.com/hibiken/asynq"
	"github.com/phongnd2802/ezy-mark/internal/pkg/email"
	"github.com/phongnd2802/ezy-mark/internal/worker"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/swagger"
	_ "github.com/phongnd2802/ezy-mark/docs"
)

// @title EzyMark Fiber API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name Philip Nguyen
// @contact.email duyphong0280@gmail.com
// @license.name MIT
// @license.url https://github.com/phongnd2802/ezy-mark/blob/main/LICENSE
// @host localhost:8000
// @BasePath /api/v1
func main() {
	app := initialize.Run()
	app.Get("/swagger/*", swagger.HandlerDefault)

	redisOpt := asynq.RedisClientOpt{
		Addr: global.Config.RedisAddr,
	}
	emailSender := email.NewGmailSender(global.Config.SenderName, global.Config.SenderEmail, global.Config.SenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, emailSender)
	
	go func() {
		log.Info().Msg("Starting Task Processor...")
		if err := taskProcessor.Start(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start task processor")
		}
	} ()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info().Msg("Starting Fiber server on :8000")
		if err := app.Listen("127.0.0.1:8000"); err != nil {
			log.Fatal().Err(err).Msg("Failed to start Fiber server")
		}
	} ()
	
	<-quit
	log.Info().Msg("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Error shutting down Fiber server")
	}

	log.Info().Msg("Server stopped")
}

