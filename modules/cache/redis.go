package cache

import (
	"github.com/gomodule/redigo/redis"
	. "mygin_websrv/conf"
	"time"
)

var RedisClient *redis.Pool

func init() {
	RedisClient = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			client, err := redis.Dial(RedisConfig["type"], RedisConfig["address"])
			if err != nil {
				return nil, err
			}
			if RedisConfig["auth"] != "" {
				if _, err := client.Do("AUTH", RedisConfig["auth"]); err != nil {
					client.Close()
					return nil, err
				}
			}
			return client, nil
		},
	}
}
