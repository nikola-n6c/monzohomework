package main

import (
	"fmt"
	"strings"
	"testing"
)

func genHTMLWithHrefs(hrefs ...string) []byte {
	var sb strings.Builder

	sb.WriteString("<html><body>")
	for _, href := range hrefs {
		sb.WriteString(fmt.Sprintf(`<a href="%s">Link</a>`, href))
	}
	sb.WriteString("</body></html>")

	return []byte(sb.String())
}

func Test_Valid_ParseStreaming(t *testing.T) {
	p := HTMLParser{}
	urls := []string{"some-url", "/someurl", "https://some.url", "https://some.url/other"}
	foundCh := make(chan string, len(urls))

	data := genHTMLWithHrefs(urls...)

	if err := p.ParseStreaming(&data, foundCh); err != nil {
		t.Error("Failed to parse", err)
	}

	for _, truth := range urls {
		url := <-foundCh
		if url != truth {
			t.Error("Incorrect parsing of href, got ", url)
		}
	}
}

func Test_ShortTag_ParseStreaming(t *testing.T) {
	p := HTMLParser{}
	foundCh := make(chan string, 1)
	truth := "some_url"

	data := []byte(fmt.Sprintf(`<html><body><a href="%s" /></body></html>`, truth))

	if err := p.ParseStreaming(&data, foundCh); err != nil {
		t.Error("Failed to parse", err)
	}

	url := <-foundCh
	if url != truth {
		t.Error("Incorrect parsing of href, got ", url)
	}
}

func Test_HrefOnOther_ParseStreaming(t *testing.T) {
	p := HTMLParser{}
	foundCh := make(chan string, 1)
	truth := "some_url"

	data := []byte(fmt.Sprintf(`<html><body><h2 href="%s" /></body></html>`, truth))

	go func() {
		if err := p.ParseStreaming(&data, foundCh); err != nil {
			t.Error("Failed to parse", err)
		}
		close(foundCh)
	}()

	url := <-foundCh
	if url == truth {
		t.Error("Incorrect parsing of href, got ", url)
	}
}

/*
// Looks like HTML parsing algorithm is much less strict than I thought
// So I'm excluding this for the sake of time
// https://godoc.org/golang.org/x/net/html#Parse
func Test_FailToParse_ParseStreaming(t *testing.T) {
	p := HTMLParser{}
	foundCh := make(chan string, 1)

	data := []byte(`nothing`)

	if err := p.ParseStreaming(&data, foundCh); err == nil {
		t.Error("Didn't fail to parse", err)
	}
}
*/
