package main

import (
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// MH stands for Monzo Homework
type MHCrawler struct {
	fetcher Fetcher
	parser  Parser
	visited StringSet
}

func (c MHCrawler) fetch(url CrawlUrl, result chan<- crawlUrlWithBody, done chan<- bool, errCh chan<- error) (isFetched bool) {
	if _, ok := c.visited[url.url.String()]; ok {
		// Already visited this URL
		done <- true
		return false
	}
	// Add to visited now to avoid computation in ShouldGet next time
	c.visited.Add(url.url.String())

	if !c.fetcher.ShouldGet(url) {
		// For whatever reason, fetcher decided we shouldn't GET this URL
		done <- true
		return false
	}

	// Start the actual fetching in another goroutine
	go func() {
		defer func() { done <- true }()

		data, err := c.fetcher.Get(url)
		if err != nil {
			errCh <- err
			return
		}

		result <- crawlUrlWithBody{&url, data}
	}()

	return true
}

func (c MHCrawler) parse(cuwb crawlUrlWithBody, results chan<- CrawlUrl, done chan<- bool, errCh chan<- error) {
	// Local channel for raw urls that parser finds
	urlsWeFound := make(chan string, 1024)

	// Start the actual parsing in another goroutine
	go func() {
		// Blocking call, the close will come after all the parsing is completed
		if err := c.parser.ParseStreaming(cuwb.body, urlsWeFound); err != nil {
			errCh <- err
		}
		close(urlsWeFound)
	}()

	// Goroutine for "stream" collection catching all the urls found by parser
	go func() {
		for foundUrl := range urlsWeFound {
			// Only look at full valid urls for now
			// maybe make this filtering more transparent and extensible?
			if parsedUrl, err := url.Parse(foundUrl); err == nil {
				results <- CrawlUrl{
					depth:  cuwb.url.depth + 1,
					parent: cuwb.url,
					url:    parsedUrl,
				}
			}
		}
		// Parser is fast so in order to see the concurrency in action
		// slow it down a bit
		// Probably the only way to debug this in a few hours
		// time.Sleep(5 * time.Second)
		done <- true
	}()
}

func (c MHCrawler) From(startUrl *url.URL) *SiteMap {
	smap := SiteMap{}

	// Concurrency control for fetchers/parsers
	// Channel capacity is arbitrary at this point
	urls := make(chan CrawlUrl, 1024)
	pageBodies := make(chan crawlUrlWithBody, 1024)

	// No need for the counters to be atomic – they're only ever modified
	// on the "main thread"
	fetchersStarted := 0
	fetcherCompleted := make(chan bool, 1024)

	parsersStarted := 0
	parserCompleted := make(chan bool, 1024)

	fetcherErrCh := make(chan error, 1)
	parserErrCh := make(chan error, 1)

	// Pop the initial url into the channel
	urls <- CrawlUrl{
		0,
		nil,
		startUrl,
	}

	// Determines wether we should stop based on all concurrency control primitives
	shouldExit := func(fs, ps int, urls chan CrawlUrl, bodies chan crawlUrlWithBody) bool {
		// No active fetchers, no active parsers, no urls and no bodies left to process
		return fs == 0 && ps == 0 && len(urls) == 0 && len(bodies) == 0
	}

forever:
	for {
		select {
		case url := <-urls:
			fetchersStarted++
			if fetched := c.fetch(url, pageBodies, fetcherCompleted, fetcherErrCh); fetched {
				smap.AddCrawlUrl(url)
			}
		case pageBody := <-pageBodies:
			parsersStarted++
			c.parse(pageBody, urls, parserCompleted, parserErrCh)
		case _ = <-fetcherCompleted:
			fetchersStarted--
			if shouldExit(fetchersStarted, parsersStarted, urls, pageBodies) {
				break forever
			}
		case _ = <-parserCompleted:
			parsersStarted--
			if shouldExit(fetchersStarted, parsersStarted, urls, pageBodies) {
				break forever
			}
		case err := <-fetcherErrCh:
			log.Error("Fetcher error: ", err)
		case err := <-parserErrCh:
			log.Error("Parser error: ", err)
		case <-time.After(5 * time.Second):
			// If the concurrency control is ok, this should not happen
			// At least not often
			// But it's good to have it to ensure that the program will end
			log.Info("Timeout, stopping Crawler")
			break forever
		}
	}

	close(urls)
	close(pageBodies)
	close(fetcherCompleted)
	close(parserCompleted)

	return &smap
}
