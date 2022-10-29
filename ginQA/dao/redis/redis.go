/*
管理redis连接和关闭
 */

package redis

import (
	"fmt"
	"ginQA/config"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)


var rdb *redis.Client


// Init 初始化连接
func Init(config *config.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		zap.L().Error("connect Redis failed", zap.Error(err))
		return
	}
	return
}

func Close() {
	_ = rdb.Close()
}