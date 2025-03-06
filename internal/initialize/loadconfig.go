package initialize

import (
	"os"

	"github.com/phongnd2802/ezy-mark/internal/config"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func loadConfig() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	if cfg.Mode == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Loading config...")
	global.Config = cfg

	log.Info().Msg("Config loaded successfully.")
}