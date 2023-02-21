package dbs

import (
	"time"

	"github.com/cheerego/go-redisson"
	"github.com/go-redis/redis/v8"
)

func Redisson(rdb *redis.Client) *godisson.Godisson {
	return godisson.NewGodisson(rdb, godisson.WithWatchDogTimeout(time.Second*30))
}

func Lock(rdb *redis.Client, key string) *godisson.Mutex {
	return Redisson(rdb).NewMutex(key)
}

func RLock(rdb *redis.Client, key string) *godisson.RLock {
	return Redisson(rdb).NewRLock(key)
}
