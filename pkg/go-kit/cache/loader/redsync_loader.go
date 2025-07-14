package cacheloader

import (
	"context"
	"fmt"
	"github.com/1612421/cinema-booking/pkg/go-kit/cache"
	"github.com/1612421/cinema-booking/pkg/go-kit/log"
	"time"

	"github.com/go-redsync/redsync/v4"
	"go.uber.org/zap"

	redislock "github.com/1612421/cinema-booking/pkg/go-kit/cache/redis/lock"
)

type RedSyncLoader[K comparable, V any] struct {
	redLock redislock.RedLock
	loader  cache.Loader[K, V]
	expiry  time.Duration
	loadKey LoadKeyFunc
}

type LoadKeyFunc func(string) string

func NewRedSyncLoader[K comparable, V any](redLock redislock.RedLock, loader cache.Loader[K, V], loadKey LoadKeyFunc, expiry time.Duration) *RedSyncLoader[K, V] {
	return &RedSyncLoader[K, V]{
		redLock: redLock,
		loader:  loader,
		expiry:  expiry,
		loadKey: loadKey,
	}
}

func (r *RedSyncLoader[K, V]) Load(ctx context.Context, c cache.Store[K, V], key K) (value V, err error) {
	mutex := r.redLock.GetLock(r.loadKey(defaultKeyEncoder(key)), redsync.WithExpiry(r.expiry))
	err = mutex.LockContext(ctx)
	if err != nil {
		return value, fmt.Errorf("acquire lock failed: %w", err)
	}

	defer func() {
		ok, errUnlock := mutex.UnlockContext(ctx)
		if !ok || errUnlock != nil {
			log.For(ctx).Error("Unlock failed", zap.Error(err))
		}
	}()

	value, err = r.loader.Load(ctx, c, key)

	return
}

func (r *RedSyncLoader[K, V]) LoadAll(ctx context.Context, c cache.Store[K, V], key K) (map[K]V, error) {
	mutex := r.redLock.GetLock(defaultKeyEncoder(key), redsync.WithExpiry(r.expiry))
	err := mutex.LockContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire lock failed: %w", err)
	}

	defer func() {
		ok, errUnlock := mutex.UnlockContext(ctx)
		if !ok || errUnlock != nil {
			log.For(ctx).Error("Unlock failed", zap.Error(err))
		}
	}()

	return r.loader.LoadAll(ctx, c, key)
}

func defaultKeyEncoder(key any) string {
	return fmt.Sprint(key)
}
