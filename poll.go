package main

import (
	"fmt"
	"net/http"
	"time"
)

func CreateGetHandler(pool *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := pool.Get()
		defer c.Close()

		output, _ := pollDataSource(w, c)

		select {
		case data := <-output:
			fmt.Fprintf(w, "{ \"data\": \"%s\" }\n", data)
		case <-time.After(time.Second):
			fmt.Fprint(w, "{}\n")
		}
	}
}
