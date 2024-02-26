package main

import (
	"fmt"
	"sync"
)

// mutex solution is very similar to counting semaphore
// limits number of goroutines acting on variable by "locking it"
// will start 1000 goroutines
func do() int {
	var m sync.Mutex // mutex
	var n int64
	var w sync.WaitGroup

	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			m.Lock()
			n++
			m.Unlock()
			w.Done()
		}()
	}

	w.Wait()      // waitGroup will make sure all 1000 goroutines will run
	return int(n) // output will be 1000 since mutex locks increment operation and makes it mutually exclusive
}

func main() {
	fmt.Println(do())
}
