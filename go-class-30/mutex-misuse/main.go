package main

import (
	"fmt"
	"sync"
	"time"
)

// NOTE: there is no way to simultaneously lock two mutexes with one command
func main() {
	var m1, m2 sync.Mutex

	done := make(chan bool)

	fmt.Println("START")

	// first goroutine locks mutex 1,
	// sleeps for one second,
	// locks mutex 2,
	// prints signal
	// unlocks mutex 2
	// unlocks mutex 1
	go func() {
		m1.Lock()
		defer m1.Unlock()
		time.Sleep(1 * time.Second)
		m2.Lock()
		defer m2.Unlock()

		fmt.Println("SIGNAL")
		done <- true
	}()

	// second goroutine, switches order of locking and unlocking m1 and m2
	// this will cause program to fail due to deadlock -> dining philosophers problem
	// YOU MUST ALWAYS ACQUIRE MUTEXES IN THE SAME PARTICULAR ORDER IF YOU ARE ACQUIRING MORE THAN ONE
	// therefore, locking and unlocking of m1 and m2 should follow same order as first goroutine
	go func() {
		m2.Lock()
		defer m2.Unlock()
		time.Sleep(1 * time.Second)
		m2.Lock()
		defer m2.Unlock()

		fmt.Println("SIGNAL")
		done <- true
	}()

	<-done
	fmt.Println("DONE")
	<-done
	fmt.Println("DONE")
}

// dining philosopher simplest scenario: there are two philosophers about to eat
// there is one knife and one fork, and you need both to eat
// one philosopher grabs a knife and the other grabs the fork
// both are waiting on the other utensil and neither of them end up eating
