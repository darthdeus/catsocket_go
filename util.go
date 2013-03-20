package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func httpError(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", text)
}

type RequestData struct {
	Valid   bool
	Data    string
	ApiKey  string
	Channel string
}

func ParseRequest(w http.ResponseWriter, req *http.Request, c Connection) RequestData {
	data := req.FormValue("data")
	apiKey := req.FormValue("api_key")
	channel := req.FormValue("channel")

	invalid := RequestData{}
	invalid.Valid = false

	if apiKey == "" || data == "" || channel == "" {
		httpError(w, 400, "channel, api_key and data are all required")
		return invalid
	}

	if c.AuthorizeKey(apiKey) {
		return RequestData{true, data, apiKey, channel}
	}

	httpError(w, 401, "authentication failed")
	return invalid
}
