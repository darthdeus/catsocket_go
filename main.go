package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

const PORT = 5000

type CatsocketHandler struct {
  Pool *DB
}

func main() {
	pool := ConnectionPool()

	fmt.Printf("Starting web server on port %d\n", PORT)

  mine := PubSubService{&pool}

  httpErr := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mine)

	check(httpErr)
}
