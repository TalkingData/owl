package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

func WriteToRedis(channel chan string, cfg *Config) {
	RedisConnPool := &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.REDIS_ADDR)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}
	for {
		rc := RedisConnPool.Get()
		select {
		case msg, ok := <-channel:
			if ok {
				_, err := rc.Do("LPUSH", cfg.REDIS_KEY, msg)
				if err != nil {
					slog.Error("rc.Do LPUSH ERROR:%v", err)
					channel <- msg
				} else {
					slog.Info("LPUSH:", msg)
				}
			}
		}
		rc.Close()
	}
}

func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		slog.Error("ping redis fail", err)
	}
	return err
}
