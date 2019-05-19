package main

// Super simple, efficient way to represend a DAG
type SiteMap map[string]StringSet

func (sm SiteMap) Add(from, to string) {
	if from != "" {
		// Add the "from" node to DAG if not there
		if _, ok := sm[from]; !ok {
			sm[from] = StringSet{}
		}
		// Add the actual link
		sm[from].Add(to)
	}

	if to != "" {
		// We're adding the "to" node to the DAG as well
		if _, ok := sm[to]; !ok {
			sm[to] = StringSet{}
		}
	}
}

func (sm SiteMap) AddCrawlUrl(crawlUrl CrawlUrl) {
	parentstr := ""
	if crawlUrl.parent != nil {
		parentstr = crawlUrl.parent.url.String()
	}

	sm.Add(parentstr, crawlUrl.url.String())
}
