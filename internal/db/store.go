package db

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface {
	Querier
}

type sqlStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &sqlStore{
		connPool: connPool,
		Queries: New(connPool),
	}
}