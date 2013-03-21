package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// Redis connection pool wrapper
type ConnectionPool struct {
	*redis.Pool
}

// Returns a single connection from the pool
func (pool ConnectionPool) Get() Connection {
	return Connection{pool.Get()}
}

func CreateConnectionPool() ConnectionPool {
	server := "0.0.0.0:6379"

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	return ConnectionPool{pool}
}
