package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

// TODO - try using a mutex to force atomic operations
type DB struct {
  pool *redis.Pool
}

func (d DB) Get() DBConnection {
  return DBConnection{d.pool.Get()}
}

type Connection interface {
  AuthorizeKey(apiKey string) bool
  Poll(channel string, score int) ([]interface{}, error)
}

type DBConnection struct {
  redis.Conn
}

func (c DBConnection) AuthorizeKey(apiKey string) bool {
	status, err := redis.Bool(c.Do("SISMEMBER", "keys", apiKey))

	if err != nil {
		panic(err)
	}

	return status
}

func (c DBConnection) Poll(channel string, score int) ([]interface{}, error) {
	reply, err := redis.Values(c.Do("ZRANGEBYSCORE", channel, 0, score))

  return reply, err
}

func ChannelName(apiKey string, name string) string {
	return apiKey + name
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

