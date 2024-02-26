package main

import "fmt"

// Prime Sieve
// numbers get generated and non-prime numbers get filtered out
// while prime numbers get added as filters
// main <----------generator // generator is func to generate numbers
// main <------------(2) generator // create goroutine that is 2-filter, take channel that goes from generator to main and make it go through 2-filter
// main <--------(3)-(2) generator // since 3 isn't filtered by 2-filter, create 3-filter
// main <--------(3)-(2) generator // 4 gets filtered out by 2-filter
// main <----(5)-(3)-(2) generator // create 5 filter, etc.
// main function will only see prime numbers and every time it sees prime number, it will create another filtering goroutine,
// hook it into the channel is gets numbers from and create new channel back to itself (main)
// https://www.youtube.com/watch?v=zJd7Dvg3XCk&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=23 

func generate(limit int, ch chan<- int) {
	for i := 2; i < limit; i++ {
		ch <- i
	}
	close(ch)
}

// read from source channel coming from generator, write to destination channel
func filter(src <-chan int, dst chan<- int, prime int) {
	for i := range src {
		if i%prime != 0 {
			dst <- i
		}
	}
	close(dst)
}

func sieve(limit int) {
	ch := make(chan int)

	go generate(100, ch)

	for {
		prime, ok := <-ch

		if !ok {
			break
		}

		// make channel
		ch1 := make(chan int)
		// filter with source being old channel and destination being new channel
		go filter(ch, ch1, prime)
		// make old channel new channel for next round in for-loop
		ch = ch1
		fmt.Print(prime, " ")
	}
}

func main() {
	// this is not really an efficient program because there is a lot of communication overhead
	// just an example for channels in goroutines
	sieve(100)
}
