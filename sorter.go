// 3 october 2014
package main

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
var toplevelID string

func assignChildren() {
	toplevels = make([]*Entry, len(entries))
	copy(toplevels, entries)
	i := 0
	for _, e := range entries {
		if parent, ok := byID[e.Parent]; ok {
			parent.Children = append(parent.Children, e)
			toplevels = append(toplevels[:i], toplevels[i + 1:]...)
		} else {
			i++
		}
	}
	for _, e := range toplevels {
		println(e.Name + " | " + e.Parent)
		println("")
	}
}
