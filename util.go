package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func pollDataSource(w http.ResponseWriter, c redis.Conn) {
	reply, err := redis.Values(c.Do("ZRANGEBYSCORE", "foo", 0, int(1e9)))

	check(err)

	if len(reply) > 0 {
		fmt.Fprintf(w, "Got a reply")
	} else {
		fmt.Fprintf(w, "Empty reply")
	}

	for _, item := range reply {
		text, _ := redis.String(item, nil)

		io.WriteString(w, string(text))
	}
}

func authorizeKey(apiKey string, c redis.Conn) bool {
	status, err := redis.Bool(c.Do("SISMEMBER", "keys", apiKey))

	if err != nil {
		panic(err)
	}

	return status
}

func ChannelName(apiKey string, name string) string {
	return apiKey + name
}

func httpError(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", text)
}

type Item struct {
	Valid   bool
	Data    string
	ApiKey  string
	Channel string
}

func ParseRequest(w http.ResponseWriter, req *http.Request, c redis.Conn) Item {
	data := req.FormValue("data")
	apiKey := req.FormValue("api_key")
	channel := req.FormValue("channel")

	invalid := Item{}
	invalid.Valid = false

	if apiKey == "" || data == "" || channel == "" {
		httpError(w, 400, "channel, api_key and data are all required")
		return invalid
	}

	// fmt.Printf("checking API key: %s\n", apiKey)

	if authorizeKey(apiKey, c) {
		return Item{true, data, apiKey, channel}
	}

	httpError(w, 401, "authentication failed")
	return invalid
}
