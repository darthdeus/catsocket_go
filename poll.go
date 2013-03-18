package main

import (
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
)

func CreateGetHandler(pool *redis.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := pool.Get()
		defer c.Close()

		go func() {
			for i := 0; i < 5; i += 1 {
				pollDataSource(w, c)
				time.Sleep(1)
			}
		}()
	}
}
