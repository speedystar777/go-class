package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// USING GOROUTINE FOR EVERY DIRECTORY AND FILE HASH

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

// on processFile is only processing a single file each time, no loop
func processFile(path string, pairs chan<- pair, wg *sync.WaitGroup, limits chan bool) {
	// second deferred function called
	defer wg.Done() // at end of function, it will signal to waitGroup that unit of work is complete

	limits <- true // next step is it attempts to get entrance into pool of active workers by pushing on limits channel

	// first deferred function called
	// we know this function call below will not block, because if we can write to limits above, then we can read from it
	defer func() { // defer requires function call, and reading from limits channel is not syntactically a function call
		<-limits // therefore, we create an anonymous function
	}() // at the end (before we signal to waitGroup), we read from limits channel to open up space

	pairs <- hashFile(path) // here we are doing actual work: hashing file and writing to pairs
}

func collectHashes(pairs <-chan pair, result chan<- results) {
	hashes := make(results)

	for p := range pairs {
		hashes[p.hash] = append(hashes[p.hash], p.path)
	}

	result <- hashes
}

func searchTree(dir string, pairs chan<- pair, wg *sync.WaitGroup, limits chan bool) error {
	defer wg.Done()

	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}

		// ignore dir itself to avoid an infinite loop!
		if fi.Mode().IsDir() && p != dir {
			wg.Add(1)
			go searchTree(p, pairs, wg, limits)
			return filepath.SkipDir
		}

		if fi.Mode().IsRegular() && fi.Size() > 0 {
			wg.Add(1)
			go processFile(p, pairs, wg, limits)
		}

		return nil
	}

	// similar sequence to processFile

	// first try to write to limits, so we can see if we are below work limit, if not we will blocked
	limits <- true

	// if we can write to limit, then we can read from it,
	// so defer this function, so that we can free up space for other workers after current worker is done
	defer func() {
		<-limits
	}()

	return filepath.Walk(dir, visit) // here we are doing the actual work
}

func run(dir string) results {
	workers := 2 * runtime.GOMAXPROCS(0)
	limits := make(chan bool, workers) // this channel will limit the number of workers/goroutines in progress (counting semaphore)
	pairs := make(chan pair, workers)  // also buffering pairs channel
	result := make(chan results)
	wg := new(sync.WaitGroup)

	// we need another goroutine so we don't block here
	go collectHashes(pairs, result)

	// multi-threaded walk of the directory tree; we need a
	// waitGroup because we don't know how many to wait for
	wg.Add(1)

	err := searchTree(dir, pairs, wg, limits)

	if err != nil {
		log.Fatal(err)
	}

	// we must close the paths channel so the workers stop
	wg.Wait()

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
