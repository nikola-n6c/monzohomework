package main

// Super quick set implementation
// struct{} because it takes 0 bytes to represent
// but it looks ugly when json-ed
type StringSet map[string]struct{}

// Dummy set element
var inset struct{}

func (ss StringSet) Add(k string) {
	ss[k] = inset
}

func (ss StringSet) IsThere(k string) bool {
	_, ok := ss[k]
	return ok
}

func (ss StringSet) AsSlice() []string {
	slice := make([]string, len(ss))
	i := 0
	for k, _ := range ss {
		slice[i] = k
		i++
	}
	return slice
}
