package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

const PORT = 5000

func main() {
	pool := CreateConnection()

	fmt.Printf("Starting web server on port %d\n", PORT)

	http.HandleFunc("/publish", CreatePublishHandler(&pool))
	http.HandleFunc("/", CreateGetHandler(&pool))

	httpErr := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	check(httpErr)
}
