package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// HTTP Handler for the main server
type CatsocketHandler struct {
	ConnectionPool *ConnectionPool
	PubSubService  PubSubService
}

// Handler constructor
func NewCatsocketHandler(pool *ConnectionPool) CatsocketHandler {
	handler := CatsocketHandler{}
	handler.ConnectionPool = pool
	handler.PubSubService = PubSubService{}

	return handler
}

func errorResponse(w http.ResponseWriter, statusCode int, text string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", text)
}

// Basic routing
func (handler CatsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := parse(r)

	connection := handler.ConnectionPool.Get()

	w.Header().Set("Content-Type", "application/json")

	if !connection.Authorize(params.apiKey) {
		errorResponse(w, 401, "Unauthorized access")
		return
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
	channel   string
	data      string
	apiKey    string
	guid      string
	timestamp int
}

func parse(r *http.Request) Params {
	params := Params{}
	params.apiKey = r.FormValue("api_key")
	params.channel = r.FormValue("channel")
	params.data = r.FormValue("data")
	params.guid = r.FormValue("guid")
	params.timestamp, _ = strconv.Atoi(r.FormValue("timestamp"))

	return params
}
