package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Problem: Want to find duplicate files based on their *content* -> use secure hash because the names/dates may differ

// for each file we have path and hash -> has represented as string instead of bytes slice so we can use it as a map key
type pair struct {
	hash, path string
}

type fileList []string
type results map[string]fileList

// below is code for hashing the files
func hashFile(path string) pair {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := md5.New() // fast and good enough, but not usually on internet

	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

// searching file path tree
func searchTree(dir string) (results, error) {
	hashes := make(results)

	err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error {
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			h := hashFile(p)
			hashes[h.hash] = append(hashes[h.hash], h.path)
		}
		return nil
	})
	return hashes, err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide dir name")
	}
	if hashes, err := searchTree(os.Args[1]); err == nil {
		for hash, files := range hashes {
			if len(files) > 1 {
				fmt.Println(hash[len(hash)-7:], len(files))

				for _, file := range files {
					fmt.Println("   ", file)
				}
			}
		}
	}
}
