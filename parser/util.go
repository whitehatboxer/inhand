package parser

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func Parse(body string) []string {
	var uris []string
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse body: %v\n", err)
	}

	for _, link := range visit(nil, doc) {
		uris = append(uris, link)
	}

	return uris
}


func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}