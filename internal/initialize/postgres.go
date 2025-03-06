package initialize

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/rs/zerolog/log"
)

func initDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connPool, err := pgxpool.New(ctx, global.Config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}

	log.Info().Msg("Pinging database...")

	err = connPool.Ping(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}

	global.ConnPool = connPool

	log.Info().Msg("Connected database successfully")
}