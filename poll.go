package main

import (
	"net/http"
	"time"
)

func CreateGetHandler(pool *DB) http.HandlerFunc {
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
