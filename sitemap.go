package main

type SiteMap struct {
	// Super simple, efficient way to represend a directed graph
	smap map[string]StringSet
	// Keep the track of root to avoid topological sort later on
	root string
}

func NewSiteMap(root string) *SiteMap {
	return &SiteMap{
		map[string]StringSet{},
		root,
	}
}

func (sm SiteMap) Has(crawlUrl CrawlUrl) bool {
	_, ok := sm.smap[crawlUrl.url.String()]
	return ok
}

func (sm SiteMap) Add(from, to string) {
	if from != "" {
		// Add the "from" node to DAG if not there
		if _, ok := sm.smap[from]; !ok {
			sm.smap[from] = StringSet{}
		}
		// Add the actual link
		sm.smap[from].Add(to)
	}

	if to != "" {
		// We're adding the "to" node to the DAG as well
		if _, ok := sm.smap[to]; !ok {
			sm.smap[to] = StringSet{}
		}
	}
}

func (sm SiteMap) Get(n string) StringSet {
	if v, ok := sm.smap[n]; ok {
		return v
	}
	return StringSet{}
}

func (sm SiteMap) AddCrawlUrl(crawlUrl CrawlUrl) {
	parentstr := ""
	if crawlUrl.parent != nil {
		parentstr = crawlUrl.parent.url.String()
	}

	sm.Add(parentstr, crawlUrl.url.String())
}
