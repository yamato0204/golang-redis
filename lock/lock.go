package main

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

var (
	ErrNotObtained = errors.New("lock: not obtained")
	ErrLockNotHeld = errors.New("lock: lock not held")
)

type Locker interface {
	TryLock(context.Context, string, time.Duration) (Lock, error)
}

type Lock interface {
	Release(context.Context) error
}

// ロック操作を提供
type locker struct {
	lockCli *redislock.Client
}

// ロックの情報保持
type lock struct {
	lock *redislock.Lock
}

func NewLocker(c *redis.Client) Locker {
	return &locker{
		lockCli: redislock.New(c),
	}
}

func NewLock(l *redislock.Lock) Lock {
	return &lock{
		lock: l,
	}
}

func (rl *locker) TryLock(ctx context.Context, key string, ttl time.Duration) (Lock, error) {
	l, err := rl.lockCli.Obtain(ctx, key, ttl, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(500*time.Millisecond), 10),
		Metadata:      "",
	})
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return nil, ErrNotObtained
		}
		return nil, err
	}
	return NewLock(l), nil
}

func (rl *lock) Release(ctx context.Context) error {
	if err := rl.lock.Release(ctx); err != nil {
		if errors.Is(err, redislock.ErrLockNotHeld) {
			return ErrLockNotHeld
		}
		return err
	}
	return nil
}
