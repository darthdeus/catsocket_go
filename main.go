package main

import "time"

type Client struct {
	in  <-chan int
	out chan int
}

func Poll(in <-chan int, out chan int) {

  select {
  case item := <-in:
    out <- item
  case <-time.After(20 * time.Millisecond):
    out <- 0
  }
}

func main() {
	data := make(chan int, 2)
	out := make(chan int)

	go Poll(data, out)

	time.Sleep(500 * time.Millisecond)

	data <- 1
	data <- 1

	println(<-out)
}
