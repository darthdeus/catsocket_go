package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Connection struct {
	redis.Conn
}

// Authorize API key against the database
func (c Connection) Authorize(apiKey string) bool {
	status, err := redis.Bool(c.Do("SISMEMBER", "keys", apiKey))

	if err != nil {
		panic(err)
	}

	return status
}

func (c Connection) Poll(channel string, score int) ([]interface{}, error) {
	reply, err := redis.Values(c.Do("ZRANGEBYSCORE", channel, score, int(1e9)))

	return reply, err
}

func (c Connection) PushData(data Params) error {
	key := time.Now().Unix()
	_, err := c.Do("ZADD", ChannelName(data.apiKey, data.channel), key, data.data)

	if err != nil {
		panic(err)
	}

	return err
}

func (c Connection) Subscribe(channelName string) (output chan []string) {
	output = make(chan []string)

	go func() {
		defer func() {
			err := c.Close()
			check(err)
		}()

		for i := 0; i < 5; i += 1 {
			response := c.pollDataSource(channelName)

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

func (c Connection) pollDataSource(channelName string) []string {
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
