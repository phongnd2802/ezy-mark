package handlers

import (
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/db"
	"github.com/phongnd2802/daily-social/internal/worker"
)

type Handler struct {
	store       db.Store
	cache       cache.Cache
	distributor worker.TaskDistributor
}

func NewHandler(store db.Store, cache cache.Cache, distributor worker.TaskDistributor) *Handler {
	return &Handler{
		store:       store,
		cache:       cache,
		distributor: distributor,
	}
}
