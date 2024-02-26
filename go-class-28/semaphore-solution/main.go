package main

import (
	"fmt"
	"sync"
)

// will start 1000 goroutines
func do() int {
	m := make(chan bool, 1) // counting semaphore of one
	var n int64
	var w sync.WaitGroup

	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			m <- true
			n++
			<-m
			w.Done()
		}()
	}

	w.Wait()      // waitGroup will make sure all 1000 goroutines will run
	return int(n) // output will be 1000 since counting semaphore limits the number of active goroutines
}

func main() {
	fmt.Println(do())
}
