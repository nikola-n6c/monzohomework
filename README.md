# Monzo homework - Web crawler

_👋 Hi there! Only on rare ocassions would I put out this much code without any documentation or writeup on what/how/why, and I felt like I should do the same for this assignment. I hope it'll give you more insight into my work on this one. Thanks! 🙏_

## Design
The entire program is made out of two high level pieces:
1. Crawler - only constructs the Site map
2. Site map renderer - takes the constructed Site map and renders it to a file

### Crawler
From the very beggining I tried to keep crawler performant and easily testable. Itself is constructed over two separate parts:
1. Fetcher - fetches raw bytes from an url
2. Parser - parses raw bytes and streams interesting pieces to a channel

Why separate the two? I believe it is:
* Somewhat easier to test them when they're apart (check out parser_test.go and fetcher_test.go)
* They can easily be implemented to support different fetching/parsing mechanisms
    * Fetcher could be implemented to crawl website from a local disk
    * Parser could be implemented to crawl raw markdown 
* More concurrent system
* Somewhat easier to test the entire Crawler once these two pieces are mocked (check out crawler_test.go)

The core Crawler itself is somewhat inspired by event driven web servers. There's a main `select` on multiple channels, and all blocking/intensive work is done on separate goroutines to keep event loop fast (crawler.go:115). Usual crawling would be highly concurrent and look something like this:
* Start url is put into thhe `urls` channel and we start the event loop
* New url _event_ happens:
    * Start a new Fetcher goroutine to get the contents
    * Return control to event loop asap and add URL to Site map if needed
* After some time, Fetcher does it's job and puts raw bytes into the `pageBody` channel
* New page body _event_ happens:
    * Start a parsing goroutine that'll put new urls into the `urls` channel
* On each Fetcher/Parser exit, we check if this is the exit condition for the crawler

This keeps repeating until there are no more urls to fetch and bytes to parse. Good thing to notice is that at any given time, we could have many Fetchers and many Parsers running. 
The output of this process is a SiteMap object (simple representation of directed graph with added utility functionalities).
### Renderer
Renderer is a super simple interface that takes in a SiteMap and spills out raw bytes into a file. It's meant to provide an extensible way to support different output formats.
I managed to provide two very basic renderers:
1. JSON renderer - useful for passing between programs (e.g. can be served with a nice frontend)
2. SVG renderer - not that useful but cool*

I've provided some SVG results in the repo itself, so you can check those out (2k x 2k px).

#### Other entities in the code
* StringSet - simple _set of strings_ implementation
* SiteMap - data structure that keeps track of directed (cyclic) graph
* interfaces.go - the high level interfaces we mentioned earlier

## Alternative designs considered
My initial idea was to make everything _streaming_. That would include incremental rendering of the SiteMap object, but I drifted away from that since it would complicate the event loop without useful benefits.
I first started with Fetcher and Parser combined, but quickly moved away because it was not really following the _separations of concernes_ and it was by definition more confusing to test.

## Improvements
If I had more time, I would:
* Invest heaviliy into testing the Crawler event loop more thoroughly. 
* Do a better job at following Golang best practices (I'm still very new to the language)
* Implement a frontend to render SiteMap in a more useful fashion

## Thanks

Thank you so much for taking the time to read this doc and the code. Have a lovely day! 🙏






