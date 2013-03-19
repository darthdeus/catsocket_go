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

func pollDataSource(w http.ResponseWriter, c DBConnection) {
  reply, err := c.Poll("foo", int(1e9))

  // TODO - this should actually just return a channel and the connection
  // should block on it until there is some data. That way we can easily
  // multiplex connections onto a few channels and re-broadcast them.
  // Maybe even return a struct containing the response and the channel?

  check(err)

  if len(reply) > 0 {
    fmt.Fprintf(w, "Got a reply\n")
  } else {
    fmt.Fprintf(w, "Empty reply\n")
  }

	for _, item := range reply {
		text, _ := redis.String(item, nil)

		io.WriteString(w, string(text))
	}
}

func httpError(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", text)
}

type Item struct {
  Valid bool
  Data string
  ApiKey string
  Channel string
}

func ParseRequest(w http.ResponseWriter, req *http.Request, c Connection) Item {
  data := req.FormValue("data")
  apiKey := req.FormValue("api_key")
  channel := req.FormValue("channel")

  invalid := Item{}
  invalid.Valid = false

  if apiKey == "" || data == "" || channel == "" {
    httpError(w, 400, "channel, api_key and data are all required")
    return invalid
  }

  if c.AuthorizeKey(apiKey) {
    return Item{true, data, apiKey, channel}
  }

  httpError(w, 401, "authentication failed")
  return invalid
}
