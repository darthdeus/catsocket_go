package main

import (
  "fmt"
  "net/http"
  "io"
  "log"
)

func Handler(w http.ResponseWriter, req *http.Request) {
  io.WriteString(w, "hello world\n")
}

const PORT = 5000

func main() {
  http.HandleFunc("/", Handler)
  err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }

  fmt.Printf("Starting web server on port %d", PORT)
}
