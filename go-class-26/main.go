package main

import (
	"fmt"
	"time"
)

type T struct {
	i byte
	b bool
}

// send function takes channel and int, create object of type T, sends pointer of T to channel
// after T object gets sent, we modify it --> IT IS DELIBERATELY A RACE CONDITION
func send(i int, ch chan<- *T) {
	t := &T{i: byte(i)}
	ch <- t
	t.b = true // UNSAFE AT ANY SPEED
}

func main() {
	vs := make([]T, 5)
	ch := make(chan *T) // if we change this to a buffer of size 5 (ch := make(chan *T, 5)) sends are non-blocking
	// BUFFERED BEHAVIOR:
	// all sends will happen and *pointers* will get put in channel
	// pointer gets modified and all go routines are done
	// after sleeping one second, we will read values from channel that have been modified

	for i := range vs {
		go send(i, ch)
	}

	time.Sleep(1 * time.Second) // all goroutines are guaranteed to have started by the time sleep is over

	// copy quickly
	for i := range vs {
		vs[i] = *<-ch // read channel and dereference
		// RENDEZVOUS BEHAVIOR:
		// ch has no buffer so we are using rendezvous model and receive/read is blocked
		// main program receives and copies (probably) before sender has a chance to return from send and modify actual value
		// basically, copying occurs before line 18 race condition
	}

	// print later to minimize time between copying and race condition
	for _, v := range vs {
		fmt.Println(v)
	}

}
