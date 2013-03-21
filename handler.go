package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// HTTP Handler for the main server
type CatsocketHandler struct {
  ConnectionPool *ConnectionPool
	PubSubService PubSubService
}

// Handler constructor
func NewCatsocketHandler(pool *ConnectionPool) CatsocketHandler {
	handler := CatsocketHandler{}
  handler.ConnectionPool = pool
	handler.PubSubService = PubSubService{}

	return handler
}

// Basic routing
func (handler CatsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := parse(r)
  connection := handler.ConnectionPool.Get()

  if !connection.Authorize(params.ApiKey) {
    fmt.Fprint(w, "Unauthorized access", 401)

  }

	if r.Method == "GET" {
		handler.PubSubService.Subscribe(w, params, connection)
	} else if r.Method == "POST" {
		handler.PubSubService.Publish(w, params, connection)
	} else {
		panic("Invalid HTTP method")
	}
}

// Params parser
type Params struct {
	Channel   string
	Data      string
	ApiKey    string
	Timestamp int
}

func parse(r *http.Request) Params {
	params := Params{}
	params.ApiKey = r.FormValue("api_key")
	params.Channel = r.FormValue("channel")
	params.Data = r.FormValue("data")
	params.Timestamp, _ = strconv.Atoi(r.FormValue("timestamp"))

	return params
}
