package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Tuple-like struct
type crawlUrlWithBody struct {
	url  *CrawlUrl
	body *[]byte
}

// This should implement the Fetcher type
// And will have the functionality of staying at one host
type StayAtHostHTTPFetcher struct {
	// Default http client in Go doesn't have a timeout
	// More here: https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	client       *http.Client
	HostToStayAt string
	MaxDepth     int
}

func NewStayAtHostHTTPFetcher(host string, depth int) *StayAtHostHTTPFetcher {
	return &StayAtHostHTTPFetcher{
		&http.Client{
			Timeout: time.Second * 30,
		},
		host,
		depth,
	}
}

func (fetcher StayAtHostHTTPFetcher) ShouldGet(crawlUrl CrawlUrl) bool {
	if crawlUrl.depth > fetcher.MaxDepth {
		// Too deep
		// log.Infof("Rejecting %s because it's too deep %d %d", crawlUrl.url.String(), crawlUrl.depth, fetcher.MaxDepth)
		return false
	}
	if crawlUrl.parent != nil && crawlUrl.parent.url == crawlUrl.url {
		// Self-anchor
		// log.Infof("Rejecting %s because it's self anchor", crawlUrl.url.String())
		return false
	}
	// log.Infof("Comparing %s %s", crawlUrl.url.Hostname(), fetcher.HostToStayAt)
	return strings.HasSuffix(crawlUrl.url.Hostname(), fetcher.HostToStayAt)
}

func (fetcher StayAtHostHTTPFetcher) Get(crawlUrl CrawlUrl) (*[]byte, error) {
	res, err := fetcher.client.Get(crawlUrl.url.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return &body, nil
}
