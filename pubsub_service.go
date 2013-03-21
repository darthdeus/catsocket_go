package main

import (
	"fmt"
	"net/http"
)

type Message struct {
	Timestamp int64    `json:"timestamp"`
	Data      []string `json:"data"`
}

// func CreatePublishHandler(pool *ConnectionPool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 		c := pool.Get()
// 		defer c.Close()
//
// 		if data := ParseRequest(w, req, c); data.Valid {
// 			c.PushData(data)
// 			fmt.Fprint(w, "OK\n")
// 		}
// 	}
// }

// func CreateGetHandler(pool *ConnectionPool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 		c := pool.Get()
//
// 		output := c.Subscribe("channel")
//
// 		data := <-output
//
// 		if data == nil {
// 			fmt.Fprint(w, "{}\n")
// 		} else {
// 			payload := Message{time.Now().Unix(), data}
//
// 			bytes, _ := json.Marshal(payload)
//
// 			fmt.Fprintf(w, "%s\n", bytes)
// 		}
// 	}
// }

type SubscriptionService interface {
	Subscribe(channel string) chan string
}

type PubSubService struct {
	*ConnectionPool
}

func (pubsub PubSubService) Subscribe(w http.ResponseWriter, params Params) {
	fmt.Fprint(w, "sub\n")
}

func valid(params Params) bool {
	return params.channel != "" &&
		params.timestamp != 0 &&
		params.apiKey != "" &&
		params.data != ""
}

func (pubsub PubSubService) Publish(w http.ResponseWriter, params Params, db Connection) {
	if !valid(params) {
		errorResponse(w, 422, "channel, data, api_key and timestamp are all required")
	} else {
		w.WriteHeader(200)
		db.PushData(params)
	}
}
