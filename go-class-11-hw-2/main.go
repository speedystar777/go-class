package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

var raw = `
<!DOCTYPE html>
<html>
  <body>
    <h1>Hello, World!</h1>
    <p>This is a paragraph.</p>
    <p>This is another paragraph.</p>
    <img src="xxx.jpg" width="104" height="142">
  </body>
</html>
`

func visit(n *html.Node, pwords, ppics *int) {
	if n.Type == html.TextNode {
		*pwords += len(strings.Fields(n.Data))
	} else if n.Type == html.ElementNode && n.Data == "img" {
		*ppics++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(c, pwords, ppics)
	}
}
func countWordsAndImages(doc *html.Node) (int, int) {
	var words, pics int
	visit(doc, &words, &pics)
	return words, pics
}

func main() {
	doc, error := html.Parse(bytes.NewReader([]byte(raw)))
	if error != nil {
		fmt.Fprintf(os.Stderr, "parse failed: %s\n", error)
	}
	words, pics := countWordsAndImages(doc)
	fmt.Printf("%d words and %d images", words, pics)

}
