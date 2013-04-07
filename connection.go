package main

import (
	"crypto/sha1"
	"fmt"
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

func (c Connection) Poll(channel string, score string) ([]interface{}, error) {
	reply, err := redis.Values(c.Do("ZRANGEBYSCORE", channel, score, int64(2e12)))

	return reply, err
}

func (c Connection) PushData(params Params) error {
	key := time.Now().UnixNano() / 1000000

	payload := fmt.Sprintf("%s|%s", params.guid, params.data)

	_, err := c.Do("ZADD", ComputeChannelName(params.apiKey, params.channel), key, payload)

	if err != nil {
		panic(err)
	}

	return err
}

func (c Connection) Subscribe(channelName string, timestamp string) (output chan []string) {
	output = make(chan []string)

	go func() {
		defer func() {
			err := c.Close()
			check(err)
		}()

		for i := 0; i < 10; i += 1 {
			response := c.pollDataSource(channelName, timestamp)

			if len(response) > 0 {
				output <- response
				return
			}

			time.Sleep(time.Millisecond * 300)
		}

		output <- nil
	}()

	return
}

func (c Connection) pollDataSource(channelName string, timestamp string) []string {
	reply, err := c.Poll(channelName, timestamp)

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

func ComputeChannelName(apiKey string, name string) string {
	hash := sha1.New()
	return fmt.Sprintf("%x", hash.Sum([]byte(apiKey+name)))
}
