package main

import (
	"bytes"
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
		params.apiKey != "" &&
		params.guid != ""
}

func (pubsub PubSubService) Subscribe(w http.ResponseWriter, params Params, db Connection) {
	if !validSubscribe(params) {
		errorResponse(w, 422, "channel, timestamp, guid and api_key are required")
		return
	}

	output := db.Subscribe(ComputeChannelName(params.apiKey, params.channel))

	messages := <-output

	if messages == nil {
		fmt.Fprint(w, "{}\n")
	} else {
		matched := []string{}

		for _, item := range messages {
			if split := bytes.Split([]byte(item), []byte("|")); string(split[0]) != params.guid {
				matched = append(matched, string(split[1]))
			}
		}

		payload := Message{time.Now().Unix(), matched}

		bytes, _ := json.Marshal(payload)

		fmt.Fprintf(w, "%s\n", bytes)
	}
}

func validPublish(params Params) bool {
	return params.channel != "" &&
		params.apiKey != "" &&
		params.data != "" &&
		params.guid != ""
}

func (pubsub PubSubService) Publish(w http.ResponseWriter, params Params, db Connection) {
	if !validPublish(params) {
		errorResponse(w, 422, "channel, data, guid and api_key are required")
		return
	}

	w.WriteHeader(200)
	db.PushData(params)
}
