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
	for _, s := range strings.Split(string(*data), ",") {
		foundCh <- s
	}
	return nil
}

func Test_From(t *testing.T) {
	// Cycles and self-refs included
	truth := SiteMap{
		smap: map[string]StringSet{
			"/": map[string]struct{}{
				"/one": inset,
				"/two": inset,
			},
			"/one": map[string]struct{}{
				"/":    inset,
				"/two": inset,
			},
			"/two": map[string]struct{}{
				"/":    inset,
				"/two": inset,
			},
		},
		root: "/",
	}
	f := DummyFetcher{
		response: map[string]string{
			"/":    `/one,/two`,
			"/one": `/,/two`,
			"/two": `/,/two`,
		},
	}

	p := DummyParser{}

	c := MHCrawler{
		f,
		p,
		StringSet{},
	}

	url, _ := url.Parse("/")
	sm := c.From(url)

	if !reflect.DeepEqual(*sm, truth) {
		t.Error("SiteMaps not equal ", fmt.Sprintf("%v\n%v", *sm, truth))
	}
}
