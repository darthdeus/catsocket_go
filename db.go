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
	reply, err := redis.Values(c.Do("ZRANGEBYSCORE", channel, score, int(1e9)))

  return reply, err
}

func (c DBConnection) PushData(data *RequestData) error {
  key := time.Now().Unix()
  _, err := c.Do("ZADD", ChannelName(data.ApiKey, data.Channel), key, data.Data)

  return err
}

func (c DBConnection) Subscribe(channelName string) (output chan []string) {
  output = make(chan []string)

  go func() {
		defer func() {
      err := c.Close()
      check(err)
    }()

    for i := 0; i < 5; i += 1 {
      response := c.pollDataSource("foo2")

      if len(response) > 0 {
        output <- response
        return
      }

      time.Sleep(time.Millisecond * 200)
    }

    output <- nil
  }()

  return
}

func (c DBConnection) pollDataSource(channelName string) []string {
	reply, err := c.Poll(channelName, 0)

	check(err)

	// TODO - this should actually just return a channel and the connection
	// should block on it until there is some data. That way we can easily
	// multiplex connections onto a few channels and re-broadcast them.
	// Maybe even return a struct containing the response and the channel?

  result := []string{}

	for _, item := range reply {
		text, _ := redis.String(item, nil)
    result = append(result, text)
	}

	return result
}

func ChannelName(apiKey string, name string) string {
	return apiKey + name
}


func ConnectionPool() DB {
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
    // TestOnBorrow: func(c redis.Conn, t time.Time) error {
    //   _, err := c.Do("PING")
    //   return err
    // },
  }

  return DB{pool}
}

