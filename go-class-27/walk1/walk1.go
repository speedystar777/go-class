package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// SEQUENTIAL FILE WALK, CONCURRENT HASHING

type pair struct {
	hash string
	path string
}

type fileList []string
type results map[string]fileList

func hashFile(path string) pair {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	hash := md5.New() // fast & good enough

	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

// the collector: gets pairs channel of input data to read from, gets results channel that it will write to
func collectHashes(pairs <-chan pair, result chan<- results) {
	hashes := make(results) // makes hash table

	for p := range pairs { // for there is no data it blocks, if there is data it reads, if it gets to close the loop is over
		hashes[p.hash] = append(hashes[p.hash], p.path) // reads pair and inserting into map until there are no more
	}

	result <- hashes // write hash table into result channel
}

// the worker
func processFiles(paths <-chan string, pairs chan<- pair, done chan<- bool) {
	for path := range paths { // one paths channel being written to by files and read by multiple workers
		pairs <- hashFile(path) // workers hash the file and write it to pairs channel until paths channel closes
	}

	done <- true // when all paths have been read, then workers are finished and we write to done channel to indicate workers are done
}

func searchTree(dir string, paths chan<- string) error {
	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}

		if fi.Mode().IsRegular() && fi.Size() > 0 {
			paths <- p
		}

		return nil
	}

	return filepath.Walk(dir, visit)
}

func run(dir string) results {
	workers := 2 * runtime.GOMAXPROCS(0) // tells us how many thread go should create based on the hardware
	// channels are buffered here
	paths := make(chan string)
	pairs := make(chan pair) // could improve time slightly by buffering pairs channel because we have a bunch of workers that are ready to report
	// results to collector, but there is only one collector, so each worker must wait until is read by collector to return (rendezvous model)
	// with buffer model, each worker can just write to pairs and return, collector will read in order that pairs are written in
	done := make(chan bool)
	result := make(chan results)

	for i := 0; i < workers; i++ {
		// as soon as these workers start, they are going to block trying to read paths because nothing has been written
		go processFiles(paths, pairs, done)
	}

	// we need another goroutine so we don't block here
	go collectHashes(pairs, result)

	if err := searchTree(dir, paths); err != nil {
		return nil
	}

	// we must close the paths channel so the workers stop
	close(paths)

	// wait for all the workers to be done
	for i := 0; i < workers; i++ {
		<-done // reading done exactly as many times as there are workers
	}

	// by closing pairs we signal that all the hashes
	// have been collected; we have to do it here AFTER
	// all the workers are done
	close(pairs)

	return <-result
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide dir name!")
	}

	if hashes := run(os.Args[1]); hashes != nil {
		for hash, files := range hashes {
			if len(files) > 1 {
				// we will use just 7 chars like git
				fmt.Println(hash[len(hash)-7:], len(files))

				for _, file := range files {
					fmt.Println("  ", file)
				}
			}
		}
	}
}
