package global

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phongnd2802/ezy-mark/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	Config   config.Config
	ConnPool *pgxpool.Pool
	Rdb      *redis.Client
)
