// 3 october 2014
package main

import (
	"sort"
)

var byID = make(map[string]*Entry)

func collectByID() {
	for _, e := range entries {
		if byID[e.ID] != nil {
			// discard older versions for now
			// TODO detect deleted content
			if byID[e.ID].Date.After(e.Date) {
				continue
			}
		}
		byID[e.ID] = e
	}
}

var toplevels []*Entry
var orphans []*Entry

func assignChildren() {
	orphans = make([]*Entry, len(entries))
	copy(orphans, entries)
	i := 0
	for _, e := range entries {
		if parent, ok := byID[e.Parent]; ok {
			parent.Children = append(parent.Children, e)
			orphans = append(orphans[:i], orphans[i + 1:]...)
		} else {
			i++
		}
	}
	i = 0
	for _, e := range orphans {
		if e.Parent == "-1" {
			toplevels = append(toplevels, e)
			orphans = append(orphans[:i], orphans[i + 1:]...)
		} else {
			i++
		}
	}
}

type sortEntries []*Entry

func (s sortEntries) Len() int { return len(s) }
func (s sortEntries) Less(i, j int) bool { return s[i].Order < s[j].Order }
func (s sortEntries) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func sortChildren() {
	for _, e := range entries {
		sort.Sort(sortEntries(e.Children))
	}
}
