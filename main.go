package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

const PORT = 5000

func main() {
	pool := CreateConnectionPool()

	fmt.Printf("Starting web server on port %d\n", PORT)

	handler := NewCatsocketHandler(&pool)

	httpErr := http.ListenAndServe(fmt.Sprintf(":%d", PORT), handler)

	check(httpErr)
}
