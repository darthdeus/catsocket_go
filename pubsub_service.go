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
	*ConnectionPool
}

func validSubscribe(params Params) bool {
	return params.channel != "" &&
		params.timestamp != 0 &&
		params.apiKey != ""
}

func (pubsub PubSubService) Subscribe(w http.ResponseWriter, params Params, db Connection) {
	if !validSubscribe(params) {
		errorResponse(w, 422, "channel, timestamp and api_key are required")
		return
	}

	output := db.Subscribe(ComputeChannelName(params.apiKey, params.channel))

	data := <-output

	if data == nil {
		fmt.Fprint(w, "{}\n")
	} else {
		payload := Message{time.Now().Unix(), data}

		bytes, _ := json.Marshal(payload)

		fmt.Fprintf(w, "%s\n", bytes)
	}
}

func validPublish(params Params) bool {
	return params.channel != "" &&
		params.apiKey != "" &&
		params.data != ""
}

func (pubsub PubSubService) Publish(w http.ResponseWriter, params Params, db Connection) {
	if !validPublish(params) {
		errorResponse(w, 422, "channel, data and api_key are required")
		return
	}

	w.WriteHeader(200)
	db.PushData(params)
}
