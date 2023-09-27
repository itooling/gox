package dbs

import (
	"fmt"
	"strings"
	"time"

	"github.com/itooling/gox/sys"
	"github.com/redis/go-redis/v9"
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
		Addr:         sys.String("redis.addr"),
		Password:     sys.String("redis.pass"),
		DB:           sys.Int("redis.index"),
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
		Addrs:        strings.Split(sys.String("redis.nodes"), ","),
		Password:     sys.String("redis.nodes_pass"),
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
