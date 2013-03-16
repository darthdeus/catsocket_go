package main

import "time"

type Client struct {
	in  <-chan int
	out chan int
}

type Terminator struct {
	done    chan bool
	blocker chan string
}

func NewTerminator() Terminator {
	return Terminator{make(chan bool), make(chan string)}
}

func Poll(in chan int, out chan int, term Terminator) {
	for {
		select {
		case item := <-in:
			out <- item
		case <-time.After(20 * time.Millisecond):
			out <- 0
		case <-term.done:
			term.blocker <- "john"
			return
		}
	}
}

func main() {
	data := make(chan int)
	out := make(chan int)

	term := NewTerminator()

	go Poll(data, out, term)

	term.done <- true

	println(<-term.blocker)

}
