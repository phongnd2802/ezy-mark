package cache

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/phongnd2802/ezy-mark/internal/global"
	"github.com/redis/go-redis/v9"
)


func GetCache(ctx context.Context, key string, obj interface{}) error {
	result, err := global.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key %s not found", key)
	} else if err != nil {
		return err
	}
	// convert result to obj
	if err := sonic.Unmarshal([]byte(result), obj); err != nil {
		return err
	}
	return nil
}