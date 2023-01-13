package dbs

import (
	"github.com/itooling/gox/tool"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
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
		Addr:         tool.String("db.redis.addr"),
		Password:     tool.String("db.redis.pass"),
		DB:           tool.Int("db.redis.index"),
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
		Addrs:        strings.Split(tool.String("db.redis.nodes"), ","),
		Password:     tool.String("db.redis.nodes_pass"),
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
