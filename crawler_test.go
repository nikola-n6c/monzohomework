package main

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type DummyFetcher struct {
	response map[string]string
}

func (df DummyFetcher) ShouldGet(crawlUrl CrawlUrl) bool {
	return true
}

func (df DummyFetcher) Get(crawlUrl CrawlUrl) (*[]byte, error) {
	val, _ := df.response[crawlUrl.url.String()]
	valbytes := []byte(val)
	return &valbytes, nil
}

type DummyParser struct{}

func (dp DummyParser) ParseStreaming(data *[]byte, foundCh chan<- string) error {
	// go func() {
	for _, s := range strings.Split(string(*data), ",") {
		foundCh <- s
	}
	// }()
	return nil
}

func Test_From(t *testing.T) {
	// Cycles included
	truth := SiteMap{
		smap: map[string]StringSet{
			"http://demo.com/": map[string]struct{}{
				"http://demo.com/one": inset,
				"http://demo.com/two": inset,
			},
			"http://demo.com/one": map[string]struct{}{
				"http://demo.com/":    inset,
				"http://demo.com/two": inset,
			},
			"http://demo.com/two": map[string]struct{}{
				"http://demo.com/": inset,
			},
		},
		root: "http://demo.com/",
	}
	f := DummyFetcher{
		response: map[string]string{
			"http://demo.com/":    `/one,/two`,
			"http://demo.com/one": `/,/two`,
			"http://demo.com/two": `/`,
		},
	}

	p := DummyParser{}

	c := MHCrawler{
		f,
		p,
		StringSet{},
	}

	url, _ := url.Parse("http://demo.com/")
	sm := c.From(url)

	if !reflect.DeepEqual(*sm, truth) {
		t.Error("SiteMaps not equal ", fmt.Sprintf("%v\n%v", *sm, truth))
	}
}
