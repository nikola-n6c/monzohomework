package main

import "net/url"

// A bit more useful than passing around raw url.URL or string
type CrawlUrl struct {
	depth  int
	parent *CrawlUrl
	url    *url.URL
}

// Fetches the content at the particular crawl url
type Fetcher interface {
	ShouldGet(CrawlUrl) bool
	Get(CrawlUrl) (*[]byte, error)
}

// Parses the particular byte slice
// Streaming here means that it'll emit parsed stuff onto a channel
type Parser interface {
	ParseStreaming(*[]byte, chan<- string) error
}
