package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// atomic solution
// making an operation atomic means that it can't be broken apart
// n++ is not an atomic operation and actually using a mutex to prevent race condition in this case is more expensive than just incrementing
// will start 1000 goroutines
func do() int {
	var n int64
	var w sync.WaitGroup

	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			atomic.AddInt64(&n, 1) // this makes the read, modify, write of n++ atomic, so it all happens together and prevents the race condition
			w.Done()
		}()
	}

	w.Wait()      // waitGroup will make sure all 1000 goroutines will run
	return int(n) // output will be 1000 since atomic operation is used, which is more primitive (done at hardware level)
}

func main() {
	fmt.Println(do())
}
