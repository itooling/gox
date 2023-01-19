package dbs

import (
	"fmt"
	"strings"
	"time"

	"gitee.com/itooling/gox"
	"github.com/go-redis/redis"
)

var rc redis.Cmdable

func Redis() redis.Cmdable {
	if rc != nil {
		return rc
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	options := redis.Options{
		Addr:         gox.String("app.dbs.redis.addr"),
		Password:     gox.String("app.dbs.redis.pass"),
		DB:           gox.Int("app.dbs.redis.index"),
		DialTimeout:  time.Millisecond * 50,
		ReadTimeout:  time.Millisecond * 50,
		WriteTimeout: time.Millisecond * 50,
	}
	rc = redis.NewClient(&options)

	return rc
}

func RedisCluster() redis.Cmdable {
	if rc != nil {
		return rc
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	options := redis.ClusterOptions{
		Addrs:        strings.Split(gox.String("app.dbs.redis.nodes"), ","),
		Password:     gox.String("app.dbs.redis.nodes_pass"),
		DialTimeout:  time.Millisecond * 50,
		ReadTimeout:  time.Millisecond * 50,
		WriteTimeout: time.Millisecond * 50,
	}
	rc = redis.NewClusterClient(&options)

	return rc
}

func RC() redis.Cmdable {
	if Redis(); rc != nil {
		return rc
	} else if RedisCluster(); rc != nil {
		return rc
	} else {
		return rc
	}
}
