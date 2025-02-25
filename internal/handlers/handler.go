package handlers

import (
	"github.com/phongnd2802/daily-social/internal/cache"
	"github.com/phongnd2802/daily-social/internal/db"
)

type Handler struct {
	store db.Store
	cache cache.Cache
}

func NewHandler(store db.Store, cache cache.Cache) *Handler {
	return &Handler{
		store: store,
		cache: cache,
	}
}