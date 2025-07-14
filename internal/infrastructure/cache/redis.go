package infrastructurecache

import (
	cacheredis "github.com/1612421/cinema-booking/pkg/go-kit/cache/redis"
	redislock "github.com/1612421/cinema-booking/pkg/go-kit/cache/redis/lock"
	"github.com/redis/go-redis/v9"

	"github.com/1612421/cinema-booking/config"
)

type RedLock redislock.RedLock

func NewRedisClient(cfg *config.Config) (redis.UniversalClient, func(), error) {
	return cacheredis.NewClient(cfg.Redis)
}

func NewRedLockClient(cfg *config.Config, redisCli redis.UniversalClient) RedLock {
	return redislock.NewRedLock(redisCli, redislock.WithPrefix(cfg.Redis.Prefix))
}
