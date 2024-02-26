package main

import (
	"log"
	"time"
)

func main() {
	chans := []chan int{
		make(chan int),
		make(chan int),
	}

	for i := range chans {
		// go routine for each channel
		go func(i int, ch chan<- int) {
			for {
				// first sleep
				time.Sleep(time.Duration(i) * time.Second)
				// return the number passed to go routine
				ch <- i
			}
		}(i+1, chans[i])
	}

	for i := 0; i < 12; i++ {
		// one of the channels is sending data twice as fast as the other --> channel 0 is sleeping for 1 second, but channel 1 is sleeping for 2
		// code below is not the best way for reading data from channels
		// we would ping-pong between these channels and not see that channel 0 is sending data every second
		// i1 := <-chans[0] // reading value from channel 0
		// i2 := <-chans[1] // until value from channel 1 is received, we will not go back and listen to channel 0

		// select will listen to both channels at same time and whichever is ready first, read from it
		// even if one channel is producing data twice as fast, we don't need to wait for other channel to produce data before we read from faster one
		// output is very different from code above; we can see that channel 0 is sending data every second
		select {
		case m0 := <-chans[0]:
			log.Println("received", m0)
		case m1 := <-chans[1]:
			log.Println("received", m1)
		}
	}
}
