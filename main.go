package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const PORT = 5000

func main() {
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
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	fmt.Printf("Starting web server on port %d\n", PORT)

	http.HandleFunc("/publish", CreatePublishHandler(pool))
	http.HandleFunc("/", CreateGetHandler(pool))

	httpErr := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	check(httpErr)
}
