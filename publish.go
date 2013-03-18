package main

import (
	"fmt"
	"net/http"
	"time"
)

func CreatePublishHandler(pool *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := pool.Get()
		defer c.Close()

		if data := ParseRequest(w, req, c); data.Valid {
			key := time.Now().Unix()

			c.Do("ZADD", ChannelName(data.ApiKey, data.Channel), key, data.Data)

			fmt.Fprint(w, "OK\n")
		}
	}
}
