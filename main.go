package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func main() {
	startUrlFlag := flag.String("starturl", "https://monzo.com", "the url to start crawling from")
	depthFlag := flag.Int("depth", 3, "the max depth to crawl to")

	flag.Parse()

	startUrl, err := url.Parse(*startUrlFlag)
	if err != nil {
		log.Fatal("Url not properly formatted")
		return
	}

	crawler := MHCrawler{
		fetcher: NewStayAtHostHTTPFetcher(startUrl.Hostname(), *depthFlag),
		parser:  HTMLParser{},
		visited: StringSet{},
	}

	siteMap := crawler.From(startUrl)

	b, err := json.MarshalIndent(siteMap, "", "  ")
	if err != nil {
		log.Error(err)
	}
	fmt.Print(string(b))
}
