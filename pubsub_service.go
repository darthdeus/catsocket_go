package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


type Message struct {
	Timestamp int64    `json:"timestamp"`
	Data      []string `json:"data"`
}

type SubscriptionService interface {
	Subscribe(channel string) chan string
}

type PubSubService struct {
  ConnectionPool *DB
}

func (pubsub PubSubService) Subscribe(w http.ResponseWriter, r *http.Request) {
}

func (pubsub PubSubService) Publish(w http.ResponseWriter, r *http.Request) {
}

func (pubsub PubSubService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
    pubsub.Subscribe(w, r)
  } else if r.Method == "POST" {
    pubsub.Publish(w, r)
  } else {
    panic("Invalid HTTP method")
  }
}

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
