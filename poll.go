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

		output := c.Subscribe("channel")

    data := <-output

    if data == nil {
      fmt.Fprint(w, "{}\n")
    } else {
      payload := Message{time.Now().Unix(), data}

      bytes, _ := json.Marshal(payload)

      fmt.Fprintf(w, "%s\n", bytes)
    }

  }
}
