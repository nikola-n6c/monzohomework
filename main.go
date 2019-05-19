package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

func renderToFile(siteMap *SiteMap, renderer SiteMapRenderer, file string) {
	buf, err := renderer.Render(siteMap)
	if err != nil {
		log.Fatal("Error rendering site map")
		return
	}

	if file == "stdout" {
		fmt.Println(buf.String())
	} else {
		f, err := os.Create(file)
		if err != nil {
			log.Fatal("Error creating output file")
			return
		}
		f.WriteString(buf.String())
		if err != nil {
			log.Fatal("Error writing output file")
			return
		}
		f.Close()
	}
}

func main() {
	startUrlFlag := flag.String("starturl", "https://monzo.com", "the url to start crawling from")
	depthFlag := flag.Int("depth", 3, "the max depth to crawl to")
	format := flag.String("format", "svg", "rendering format of the sitemap (json, svg)")
	file := flag.String("file", "stdout", "file to write to (omit for stdout)")

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

	var renderer SiteMapRenderer
	switch *format {
	case "svg":
		renderer = SVGSiteMapRenderer{}
	case "json":
		renderer = JSONSiteMapRenderer{}
	}

	renderToFile(siteMap, renderer, *file)
}
