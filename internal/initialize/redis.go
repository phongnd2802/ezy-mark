package initialize

import (
	"context"
	"time"

	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func initRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     global.Config.RedisAddr,
		Password: "secret",
		DB:       0,
		PoolSize: 10,
	})


	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	log.Info().Msg("Pinging redis...")

	err := client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	global.Rdb = client

	log.Info().Msg("Connected to redis successfully")
}