package redis

import (
	"github.com/1612421/cinema-booking/config"
	"github.com/1612421/cinema-booking/internal/entity"
	"github.com/1612421/cinema-booking/pkg/go-kit/cache"
	cacheredis "github.com/1612421/cinema-booking/pkg/go-kit/cache/redis"
	timesdk "github.com/1612421/cinema-booking/pkg/go-kit/time"
	"github.com/redis/go-redis/v9"
)

type IBookingCache interface {
}

type BookingCache struct {
	cache cache.Store[string, *entity.Booking]
}

const BookingCacheKey = "booking"

var BookingTTL = timesdk.MinuteDuration(10)

func NewBookingCache(cfg *config.Config, redisCli redis.UniversalClient) *BookingCache {
	return &BookingCache{
		cache: cacheredis.NewRedisCache[string, *entity.Booking](
			redisCli,
			cacheredis.WithPrefix[string, *entity.Booking](cfg.Redis.Prefix),
		),
	}
}
