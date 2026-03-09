package redis

import (
	"github.com/go-redis/redis"
	"log"
)

var RedisClient *redis.Client

func InitRedis(addr string, password string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 没有设置密码
		DB:       0,        // 使用默认DB
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("连接Redis失败：%v", err)
	}
}
