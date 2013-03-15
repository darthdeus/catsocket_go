package main

import "time"

func Poll(inbox <-chan int, out chan int) {
  for item := range inbox {
    println("Got", item, "sending back to client socket")
    out <- item
    return
  }
}

func main() {
  dataSource := make(chan int)
  clientSocket := make(chan int)

  go Poll(dataSource, clientSocket)

  time.Sleep(500 * time.Millisecond)

  dataSource <- 1

  <-clientSocket
}
