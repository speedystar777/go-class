package main

import (
	"sync"
	"testing"
)

// go test -bench=. -benchtime=10s -cpu=2,4,8 ./false_sharing_test.go
// runs for 10 second with 2, 4, of 8 threads

const (
	nchan   = 8
	nworker = nchan
	buffer  = 1024
)

var wg sync.WaitGroup

func count_share(cnt *uint64, in <-chan int) {
	for i := range in {
		*cnt += uint64(i) // false sharing, much less efficient
	}

	wg.Done()
}

func count_dont_share(cnt *uint64, in <-chan int) {
	var total int

	for i := range in {
		total += i
	}

	*cnt = uint64(total)
	wg.Done()
}

func fill(n int, in chan<- int) {
	for i := 0; i < n; i++ {
		in <- i
	}
	close(in)
}

func run(count func(cnt *uint64, int <-chan int)) (total int) {
	var cnt [nworker]uint64

	in := make([]chan int, nchan)

	for i := 0; i < nchan; i++ {
		in[i] = make(chan int, buffer)
		go fill(10000, in[i])
	}

	wg.Add(nworker)

	for i := 0; i < nworker; i++ {
		go count(&cnt[i], in[i%nchan])
	}

	wg.Wait()

	for _, v := range cnt {
		total += int(v)
	}
	return
}

func BenchmarkShare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(count_share)
	}
}

func BenchmarkDontShare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(count_dont_share)
	}
}
