package main

import (
	"fmt"
	"net/http"
)

func CreatePublishHandler(pool *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := pool.Get()
		defer c.Close()

		if data := ParseRequest(w, req, c); data.Valid {
			c.PushData(&data)
			fmt.Fprint(w, "OK\n")
		}
	}
}
