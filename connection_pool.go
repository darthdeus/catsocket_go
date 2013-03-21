package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// Redis connection pool wrapper
type ConnectionPool struct {
	RedisPool *redis.Pool
}

// Returns a single connection from the pool
func (pool ConnectionPool) Get() Connection {
	return Connection{pool.RedisPool.Get()}
}

func CreateConnectionPool() ConnectionPool {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
      c, err := redis.Dial("tcp", "0.0.0.0:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	return ConnectionPool{pool}
}
