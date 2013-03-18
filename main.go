package main

import (
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

func CreateGetHandler(in chan string, c redis.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		for i := 0; i < 5; i += 1 {
			pollDataSource(w, c)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func PushData(in chan<- string) {
	for {
		time.Sleep(time.Second)
		in <- "randomness"
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const PORT = 5000

func main() {
	c, err := redis.Dial("tcp", ":6379")
	check(err)

	in := make(chan string)

	go PushData(in)

	handler := CreateGetHandler(in, c)

	fmt.Printf("Starting web server on port %d", PORT)

	http.HandleFunc("/", handler)
	httpErr := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	check(httpErr)
}
