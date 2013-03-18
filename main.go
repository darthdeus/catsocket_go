package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func pollDataSource(w http.ResponseWriter, c redis.Conn) {
	reply, _ := redis.Values(c.Do("ZRANGEBYSCORE", "foo", 0, int(1e9)))

	for _, item := range reply {
		// text := string(item.([]byte))
		text, _ := redis.String(item, nil)
		fmt.Printf("%s", text)
		io.WriteString(w, string(text))
	}
}

func createGetHandler(c redis.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "GET" {
			w.WriteHeader(400)
			io.WriteString(w, "{ \"error\": \"Poll only supports GET\" }")
			return
		}

		go func() {
			for i := 0; i < 5; i += 1 {
				pollDataSource(w, c)
				time.Sleep(1)
			}
		}()
	}
}

func authorizeKey(apiKey string, c redis.Conn) bool {
	status, _ := redis.Bool(c.Do("SISMEMBER", "keys", apiKey))
	return status
}

func createPublishHandler(c redis.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
    req.ParseForm()

    apiKey := req.Form["api_key"][0]

    fmt.Printf("checking API key: %s\n", apiKey)

    if authorizeKey(apiKey, c) {
      io.WriteString(w, "OK\n")
    } else {
      w.WriteHeader(401)
      io.WriteString(w, "ERROR\n")
    }
	}
}

const PORT = 5000

func main() {
	c, redisErr := redis.Dial("tcp", ":6379")
	check(redisErr)

	fmt.Printf("Starting web server on port %d\n", PORT)

	http.HandleFunc("/", createGetHandler(c))
	http.HandleFunc("/publish", createPublishHandler(c))

	httpErr := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	check(httpErr)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
