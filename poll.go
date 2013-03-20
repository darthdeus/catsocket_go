package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Message struct {
  Timestamp int64 `json:"timestamp"`
  Data []string `json:"data"`
}

func CreateGetHandler(pool *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		c := pool.Get()
		defer c.Close()

		output := c.Subscribe("channel")

		select {
		case data := <-output:
      payload := Message{time.Now().Unix(), data}

      bytes, _ := json.Marshal(payload)

			fmt.Fprintf(w, "%s\n", bytes)
		case <-time.After(time.Second):
			fmt.Fprint(w, "{}\n")
		}
	}
}
