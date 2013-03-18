package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type DB struct {
  *redis.Pool
}

func CreateConnection() DB {
  server := "0.0.0.0:6379"

  pool := &redis.Pool{
    MaxIdle: 3,
    IdleTimeout: 240 * time.Second,
    Dial: func () (redis.Conn, error) {
      c, err := redis.Dial("tcp", server)
      if err != nil {
        return nil, err
      }
      return c, err
    },
    TestOnBorrow: func(c redis.Conn, t time.Time) error {
      _, err := c.Do("PING")
      return err
    },
  }

  return DB{pool}
}

