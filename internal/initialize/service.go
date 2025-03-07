package initialize

import (
	"github.com/hibiken/asynq"
	"github.com/phongnd2802/ezy-mark/internal/cache"
	"github.com/phongnd2802/ezy-mark/internal/db"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/phongnd2802/ezy-mark/internal/services"
	"github.com/phongnd2802/ezy-mark/internal/services/impl"
	"github.com/phongnd2802/ezy-mark/internal/worker"
)

func initServiceInterfaces() {
	store := db.NewStore(global.ConnPool)
	cache := cache.NewRedisClient(global.Rdb)

	// Distributor
	redisOpt := asynq.RedisClientOpt{
		Addr: global.Config.RedisAddr,
	}
	distributor := worker.NewRedisTaskDistributor(redisOpt)

	services.InitAuthService(impl.NewAuthServiceImpl(store, cache, distributor))
	services.InitUserInfo(impl.NewUserInfoImpl(store, cache))
}
