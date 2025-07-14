package redislock

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -destination=./mocks/$GOFILE -source=$GOFILE -package=redislock
type LockMutex interface {
	// TryLockContext only attempts to lock m once and returns immediately regardless of success or failure without retrying.
	TryLockContext(ctx context.Context) error
	// LockContext locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
	LockContext(ctx context.Context) error
	// UnlockContext unlocks m and returns the status of unlock.
	UnlockContext(ctx context.Context) (bool, error)
}

type RedLock interface {
	GetLock(key string, options ...redsync.Option) LockMutex
}

type redisRedLock struct {
	opts      *Options
	redisLock *redsync.Redsync
}

func NewRedLock(redisCli redis.UniversalClient, options ...Option) RedLock {
	pool := goredis.NewPool(redisCli)
	redisLock := redsync.New(pool)

	opts := newDefaultOption()
	for _, o := range options {
		o(opts)
	}
	return &redisRedLock{
		opts:      opts,
		redisLock: redisLock,
	}
}

func (r *redisRedLock) GetLock(key string, options ...redsync.Option) LockMutex {
	return r.redisLock.NewMutex(r.buildKey(key), options...)
}

func (r *redisRedLock) buildKey(key string) string {
	return r.opts.prefix + ":" + key + ":" + r.opts.suffix
}
