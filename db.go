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

func (c DBConnection) PushData(data *RequestData) error {
  key := time.Now().Unix()
  _, err := c.Do("ZADD", ChannelName(data.ApiKey, data.Channel), key, data.Data)

  return err
}

func (c DBConnection) Subscribe(chanel string) {
}

func pollDataSource(w http.ResponseWriter, c DBConnection) (output chan string, timeout chan bool) {
	reply, err := c.Poll("foo", int(1e9))

	// TODO - this should actually just return a channel and the connection
	// should block on it until there is some data. That way we can easily
	// multiplex connections onto a few channels and re-broadcast them.
	// Maybe even return a struct containing the response and the channel?

	output = make(chan string)
	timeout = make(chan bool)

	go func() {
		// TODO - the actual implementation goes here
		time.Sleep(time.Millisecond * 20)
		output <- "john"
	}()

	check(err)

  // TODO - buffer this instead and send it to the channel at once
	for _, item := range reply {
		text, _ := redis.String(item, nil)

		io.WriteString(w, string(text))
	}

	return
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

