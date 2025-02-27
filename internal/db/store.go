package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phongnd2802/daily-social/internal/global"
)

type Store interface {
	Querier
}

type sqlStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore() Store {
	return &sqlStore{
		connPool: global.ConnPool,
		Queries: New(global.ConnPool),
	}
}