package main

import (
  _ "net/http/pprof"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"net/http"
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

    go func() {
      for i := 0; i < 5; i += 1 {
        pollDataSource(w, c)
        time.Sleep(1 * time.Millisecond)
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
    io.WriteString(w, "OK")
  }
}

const PORT = 5000

func main() {
	c, err := redis.Dial("tcp", ":6379")
	check(err)

  println(authorizeKey("foo", c))

	fmt.Printf("Starting web server on port %d", PORT)

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

