package main

import (
  "fmt"
  "net/http"
  "io"
)

func Handler(w http.ResponseWriter, req *http.Request) {
  io.WriteString(w, "hello world\n")
}

const PORT = 5000

func main() {
  http.HandleFunc("/", Handler)
  http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

  fmt.Printf("Starting web server on port %d", PORT)
}
