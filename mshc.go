// 3 october 2014

package main

import (
	"fmt"
	"os"
	"io"
	"archive/zip"
	"path/filepath"
	"code.google.com/p/go.net/html"
	"encoding/json"
)

type Entry struct {
	Name	string
	Book		string
	ID		string
	Parent	string
	Order	string		// zero-based
	File		string
}

var entries []Entry

func parseEntry(r io.Reader, filename string) {
	var e Entry

	e.File = filename
	t := html.NewTokenizer(r)
	for {
		tt := t.Next()
		if tt == html.ErrorToken {
			err := t.Err()
			if err == io.EOF {
				break
			}
			panic(err)		// TODO
		}
		tok := t.Token()
		switch tok.Type {
		case html.StartTagToken, html.SelfClosingTagToken:
			if tok.Data != "meta" {
				break
			}

			var where *string
			var what string

			for _, a := range tok.Attr {
				if a.Key == "name" {
					switch a.Val {
					case "Title":
						where = &e.Name
					case "Microsoft.Help.Book":
						where = &e.Book
					case "Microsoft.Help.Id":
						where = &e.ID
					case "Microsoft.Help.TocParent":
						where = &e.Parent
					case "Microsoft.Help.TocOrder":
						where = &e.Order
					}
				} else if a.Key == "content" {
					what = a.Val
				}
			}
			if where != nil {
				*where = what
			}
		}
	}
	entries = append(entries, e)
}

func main() {
	z, err := zip.OpenReader(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer z.Close()
	for _, f := range z.File {
		if e := filepath.Ext(f.Name); e != ".htm" {
			continue
		}
		r, err := f.Open()
		if err != nil {
			panic(err)		// TODO
		}
		parseEntry(r, f.Name)
		r.Close()
	}
	b, err := json.MarshalIndent(entries, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", b)
}
