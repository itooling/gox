package dbs

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	lockScript = `
	if redis.call('exists',KEYS[1]) == 0 then
		redis.call('set',KEYS[1], 1)
		redis.call('pexpire',KEYS[1],ARGV[1])
		return 0
	end
	return redis.call('pttl',KEYS[1])`

	unlockScript = `
	if (redis.call('exists', KEYS[1]) == 0) then
		return 0;
	end;
	redis.call('del', KEYS[1]);
	return 1;`

	defaultLeaseDuration time.Duration = 30 * time.Second
)

type Mutex struct {
	client  redis.UniversalClient
	key     string
	options Options
}

type Options struct {
	leaseDuration time.Duration
}

func defaultOptions() *Options {
	opts := &Options{}
	WithLeaseDuration(defaultLeaseDuration)(opts)
	return opts
}

type Option func(*Options)

func WithLeaseDuration(leaseDuration time.Duration) Option {
	return func(mo *Options) {
		mo.leaseDuration = leaseDuration
	}
}

func NewMutex(client redis.UniversalClient, key string, options ...Option) *Mutex {
	opts := defaultOptions()
	for _, option := range options {
		option(opts)
	}
	return &Mutex{key: key, client: client, options: *opts}
}

func (m *Mutex) TryLock(ctx context.Context) (bool, error) {
	ttl, err := m.doTryLock(ctx)
	if err != nil {
		return false, fmt.Errorf("trylock failed: %w", err)
	}

	return ttl == 0, nil
}

func (m *Mutex) doTryLock(ctx context.Context) (int64, error) {
	ttl, err := m.client.Eval(ctx, lockScript, []string{m.key}, m.options.leaseDuration.Milliseconds()).Result()
	if err != nil {
		return 0, err
	}

	return ttl.(int64), nil
}

func (m *Mutex) Unlock(ctx context.Context) error {
	_, err := m.client.Eval(ctx, unlockScript, []string{m.key}, m.options.leaseDuration.Milliseconds()).Int64()
	if err != nil {
		return fmt.Errorf("unlocked failed: %w", err)
	}

	return nil
}

type Redisson struct {
	redisClient *redis.Client
}

func NewRedisson(redisClient *redis.Client) *Redisson {
	return &Redisson{redisClient: redisClient}
}

func (r *Redisson) NewMutex(key string, options ...Option) *Mutex {
	return NewMutex(r.redisClient, key, options...)
}
