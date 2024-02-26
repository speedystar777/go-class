package main

import (
	"log"
	"time"
)

func main() {
	// one tick every 2 seconds
	const tickRate = 2 * time.Second
	// once we reach five ticks worth of time, channel will be ready to read
	stopper := time.After(5 * tickRate)
	// NewTicker returns ticker object and .C gets channel of ticker object
	ticker := time.NewTicker(tickRate).C

	log.Println("start")

loop:
	for {
		select {
		// every 2 seconds print tick
		case <-ticker:
			log.Println("tick")
		// stop after five ticks
		case <-stopper:
			break loop
		}
	}

	log.Println("finish")
}
