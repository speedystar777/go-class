package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func load(index int) []byte {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", index)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "skipping %d: got %d\n", index, resp.StatusCode)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return data
}

func main() {
	var (
		output      io.WriteCloser = os.Stdout
		err         error
		total_count int
		count_404   int
		data        []byte
	)

	if len(os.Args) > 1 {
		output, err = os.Create(os.Args[1])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		defer output.Close()
	}

	fmt.Fprint(output, "[")
	defer fmt.Fprint(output, "]")

	for index := 1; count_404 < 2; index++ {
		if data = load(index); data == nil {
			count_404++
			continue
		}

		if total_count > 0 {
			fmt.Fprint(output, ",")
		}

		_, err := io.Copy(output, bytes.NewBuffer(data))

		if err != nil {
			fmt.Fprintf(os.Stderr, "stopped: %s\n", err)
			os.Exit(-1)
		}

		count_404 = 0
		total_count++
	}

	fmt.Fprintf(os.Stderr, "read %d comics\n", total_count)
}
