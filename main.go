package main

import (
	"fmt"
	"github.com/gosexy/redis"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func pollDataSource(w http.ResponseWriter, c *redis.Client) {
	reply, _ := c.ZRangeByScore("foo", 0, int(1e9))

	for _, item := range reply {
		fmt.Printf("%s", item)
		io.WriteString(w, item)
	}
}

func createGetHandler(c *redis.Client) http.HandlerFunc {
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

func authorizeKey(apiKey string, c *redis.Client) bool {
	status, err := c.SIsMember("keys", apiKey)

	if err != nil {
		panic(err)
	}

	return status
}

func channelName(apiKey string, name string) string {
	return apiKey + name
}

func httpError(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{ \"error\": \"%s\" }\n", text)
}

func createPublishHandler(c *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := req.FormValue("data")
		apiKey := req.FormValue("api_key")
		channel := req.FormValue("channel")

		if apiKey == "" || data == "" || channel == "" {
			httpError(w, 400, "channel, api_key and data are all required")
			return
		}

		fmt.Printf("checking API key: %s\n", apiKey)

		if authorizeKey(apiKey, c) {
			key := time.Now().Unix()

			status, err := c.ZAdd(channelName(apiKey, channel), key, data)
			println(status)

			if err != nil {
				panic(err)
			}

			println("added")
			io.WriteString(w, "OK\n")
		} else {
			httpError(w, 401, "authentication failed")
		}
	}
}

const PORT = 5000

func main() {
	c := redis.New()
	redisErr := c.Connect("0.0.0.0", 6379)

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
