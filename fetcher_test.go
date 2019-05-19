package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Valid_ShouldGet(t *testing.T) {
	depth := 1
	host := "demo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("http://demo.com/example")
	crawlUrl := CrawlUrl{
		depth:  0,
		parent: nil,
		url:    url,
	}

	if !f.ShouldGet(crawlUrl) {
		t.Error("Should get ", crawlUrl.url.String())
	}
}

func Test_Subdomain_ShouldGet(t *testing.T) {
	depth := 1
	host := "demo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("http://subdomain.demo.com")
	crawlUrl := CrawlUrl{
		depth:  1,
		parent: nil,
		url:    url,
	}

	if !f.ShouldGet(crawlUrl) {
		t.Error("Should get ", crawlUrl.url.String())
	}
}

func Test_Scheme_ShouldGet(t *testing.T) {
	depth := 1
	host := "demo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("https://demo.com")
	crawlUrl := CrawlUrl{
		depth:  1,
		parent: nil,
		url:    url,
	}

	if !f.ShouldGet(crawlUrl) {
		t.Error("Should get ", crawlUrl.url.String())
	}
}

func Test_WithQueryParams_ShouldGet(t *testing.T) {
	depth := 1
	host := "demo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("http://demo.com/demo?foo=bar")
	crawlUrl := CrawlUrl{
		depth:  1,
		parent: nil,
		url:    url,
	}

	if !f.ShouldGet(crawlUrl) {
		t.Error("Should get ", crawlUrl.url.String())
	}
}

func Test_TooDeep_ShouldGet(t *testing.T) {
	depth := 2
	host := "demo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("http://demo.com/demo")
	crawlUrl := CrawlUrl{
		depth:  3,
		parent: nil,
		url:    url,
	}

	if f.ShouldGet(crawlUrl) {
		t.Error("Should not get ", crawlUrl.url.String())
	}
}

func Test_TooDeep2_ShouldGet(t *testing.T) {
	depth := 2
	host := "demo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("http://demo.com/demo")
	crawlUrl := CrawlUrl{
		depth:  3,
		parent: nil,
		url:    url,
	}

	if f.ShouldGet(crawlUrl) {
		t.Error("Should not get ", crawlUrl.url.String())
	}
}

func Test_SelfRef_ShouldGet(t *testing.T) {
	depth := 1
	host := "demo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("http://demo.com/demo")
	parent := CrawlUrl{
		depth:  1,
		parent: nil,
		url:    url,
	}
	crawlUrl := CrawlUrl{
		depth:  1,
		parent: &parent,
		url:    url,
	}

	if f.ShouldGet(crawlUrl) {
		t.Error("Should not get ", crawlUrl.url.String())
	}
}

func Test_WrongHost_ShouldGet(t *testing.T) {
	depth := 1
	host := "anotherdemo.com"
	f := NewStayAtHostHTTPFetcher(host, depth)

	url, _ := url.Parse("http://demo.com/demo")
	crawlUrl := CrawlUrl{
		depth:  1,
		parent: nil,
		url:    url,
	}

	if f.ShouldGet(crawlUrl) {
		t.Error("Should not get ", crawlUrl.url.String())
	}
}

func Test_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`content`))
	}))
	defer server.Close()

	depth := 1
	host := server.URL
	f := StayAtHostHTTPFetcher{
		client:       server.Client(),
		HostToStayAt: host,
		MaxDepth:     depth,
	}

	url, _ := url.Parse(server.URL)
	fetchUrl := CrawlUrl{
		depth:  0,
		parent: nil,
		url:    url,
	}

	data, err := f.Get(fetchUrl)
	if err != nil {
		t.Error("Error fetching ", err)
	}

	if bytes.Compare(*data, []byte("content")) != 0 {
		t.Error("Error, response doesn't match")
	}
}
