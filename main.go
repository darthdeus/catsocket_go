package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func CreateGetHandler(in chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, <-in)
	}
}

func PushData(in chan<- string) {
	for {
		time.Sleep(time.Second)
		in <- "randomness"
	}
}

const PORT = 5000

func main() {
	in := make(chan string)

	go PushData(in)

	handler := CreateGetHandler(in)

	http.HandleFunc("/", handler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	fmt.Printf("Starting web server on port %d", PORT)
}
