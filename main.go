package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"net/http"
	"time"
)

func CreateGetHandler(in chan string, c *redis.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reply, _ := redis.String((*c).Do("GET", "foo"))

		<-in
		io.WriteString(w, reply)
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

	handler := CreateGetHandler(in, &c)

	http.HandleFunc("/", handler)
	httpErr := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	check(httpErr)

	fmt.Printf("Starting web server on port %d", PORT)
}
