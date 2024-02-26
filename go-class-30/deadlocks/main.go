package main

import "fmt"

// this program will fail because it has a deadlock
func main() {
	ch := make(chan bool) // channel is not buffered so send and receive are synchronized
	// i.e., if nothing is sent/written, nothing is received/read

	// starting a goroutine that may or may not put something on the channel
	go func(ok bool) {
		fmt.Println("START")
		if ok {
			ch <- ok
		}
	}(false) // since we call goroutine with false, we will not put anything on channel

	<-ch // then we are waiting to read channel, which may or may not have a value on it
	// overall, program will never move forward if we don't put anything on channel
	fmt.Println("DONE")
}
