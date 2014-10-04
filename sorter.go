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
