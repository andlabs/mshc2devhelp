// 3 october 2014
package main

import (
	"fmt"
)

var byID = make(map[string]*Entry)

func collectByID() {
	for _, e := range entries {
		if byID[e.ID] != nil {
			panic(fmt.Sprintf("duplicate %q: %#v vs %#v", e.ID, byID[e.ID], e))
		}
		byID[e.ID] = e
	}
}
