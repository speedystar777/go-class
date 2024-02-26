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

// CONCURRENT FILE WALK, CONCURRENT HASHING

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

func processFiles(paths <-chan string, pairs chan<- pair, done chan<- bool) {
	for path := range paths {
		pairs <- hashFile(path)
	}

	done <- true
}

func collectHashes(pairs <-chan pair, result chan<- results) {
	hashes := make(results)

	for p := range pairs {
		hashes[p.hash] = append(hashes[p.hash], p.path)
	}

	result <- hashes
}

func searchTree(dir string, paths chan<- string, wg *sync.WaitGroup) error {
	defer wg.Done() // after we finish searching this subtree, we remove one count from the wait group

	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}

		// ignore dir itself to avoid an infinite loop!
		if fi.Mode().IsDir() && p != dir {
			wg.Add(1) // adding before we start a new tree search
			go searchTree(p, paths, wg)
			return filepath.SkipDir
		}

		if fi.Mode().IsRegular() && fi.Size() > 0 {
			paths <- p
		}

		return nil
	}

	return filepath.Walk(dir, visit)
}

func run(dir string) results {
	workers := 2 * runtime.GOMAXPROCS(0)
	paths := make(chan string)
	pairs := make(chan pair)
	done := make(chan bool)
	result := make(chan results)
	wg := new(sync.WaitGroup) // when trying to walk the file tree in parallel, we don't know in advance how many goroutines to start
	// therefore, we can't wait for workers to be done, by iterating through number of workers and reading done channelxs
	// wait group is a tracker that can be used for this situation

	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}

	// we need another goroutine so we don't block here
	go collectHashes(pairs, result)

	// multi-threaded walk of the directory tree; we need a
	// waitGroup because we don't know how many to wait for
	wg.Add(1) // waitGroup has baseline of zero, every time we start searching down a tree, we add one to the waitGroup
	// when the work is finished, the count goes down, eventually it will go back down to zero
	// and when it goes to zero, we have finished all our work

	err := searchTree(dir, paths, wg)

	if err != nil {
		log.Fatal(err)
	}

	// we must close the paths channel so the workers stop
	wg.Wait() // this is tracking to see when waitGroup count goes back down to zero, we indicates entire file tree is searched
	close(paths)

	// wait for all the workers to be done
	for i := 0; i < workers; i++ {
		<-done
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
