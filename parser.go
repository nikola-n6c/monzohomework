package main

import (
	"bytes"

	"golang.org/x/net/html"
)

type HTMLParser struct{}

func (parser HTMLParser) ParseStreaming(data *[]byte, foundCh chan<- string) error {
	// Taken straight from the docs
	// https://godoc.org/golang.org/x/net/html

	r := bytes.NewReader(*data)
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					foundCh <- attr.Val
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return nil
}
