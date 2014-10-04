// 3 october 2014

package main

import (
	"os"
	"io"
	"archive/zip"
	"path/filepath"
	"strconv"
	"time"
	"unsafe"
	"io/ioutil"
	"code.google.com/p/go.net/html"
)

// #cgo LDFLAGS: -lmspack
// #include <mspack.h>
// #include <stdlib.h>
// /* because cgo can't handle function pointers */
// struct mscabd_cabinet *cabOpen(struct mscab_decompressor *c, char *filename)
// {
// 	return (*c->open)(c, filename);
// }
// void cabClose(struct mscab_decompressor *c, struct mscabd_cabinet *cab)
// {
// 	(*c->close)(c, cab);
// }
// int cabExtract(struct mscab_decompressor *c, struct mscabd_file *in, char *out)
// {
// 	return (*c->extract)(c, in, out);
// }
// int cabLastError(struct mscab_decompressor *c)
// {
// 	return (*c->last_error)(c);
// }
import "C"

type Entry struct {
	Name	string
	Book		string
	ID		string
	Parent	string
	Order	int
	Date		time.Time
	MSHC	string
	File		string

	Children	[]*Entry
}

var entries []*Entry

func parseEntry(r io.Reader, mshcname string, filename string) {
	var e *Entry

	e = new(Entry)
	e.MSHC = mshcname
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
			var order string
			var date string
			var err error

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
						where = &order
					case "Microsoft.Help.TopicPublishDate":
						where = &date
					}
				} else if a.Key == "content" {
					what = a.Val
				}
			}
			if where != nil {
				*where = what
			}
			if where == &order {
				e.Order, err = strconv.Atoi(order)
				if err != nil {
					panic(err)		// TODO
				}
			} else if where == &date {
				e.Date, err = time.Parse(time.RFC1123, date)
				if err != nil {
					panic(err)		// TODO
				}
			}
		}
	}
	entries = append(entries, e)
}

func parseMSHC(mshcname string) {
	z, err := zip.OpenReader(mshcname)
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
		parseEntry(r, mshcname, f.Name)
		r.Close()
	}
}

func parseCAB(cabname string, workdir string) {
	c := C.mspack_create_cab_decompressor(nil)
	if c == nil {
		panic("mspack_create_cab_decompressor() failed (TOOD error)")
	}
	defer C.mspack_destroy_cab_decompressor(c)
	ccabname := C.CString(cabname)
	defer C.free(unsafe.Pointer(ccabname))
	cab := C.cabOpen(c, ccabname)
	if cab == nil {
		panic("error opening cabinet file (TODO get error)")
	}
	defer C.cabClose(c, cab)
	for cf := cab.files; cf != nil; cf = cf.next {
		filename := C.GoString(cf.filename)
		if e := filepath.Ext(filename); e == ".mshc" {
			tmpfile := filepath.Join(workdir, filepath.Base(filename))
			ctmpfile := C.CString(tmpfile)
			err := C.cabExtract(c, cf, ctmpfile)
			if err != C.MSPACK_ERR_OK {
				panic("error extracting (TODO error to string)")
			}
			C.free(unsafe.Pointer(ctmpfile))
			parseMSHC(tmpfile)
		}
	}
}

func main() {
	workdir, err := ioutil.TempDir("", "mshc")
	if err != nil {
		panic(err)		// TODO
	}
	for _, cab := range os.Args[1:] {
		parseCAB(cab, workdir)
	}
	collectByID()
	assignChildren()
	sortChildren()
	buildDevhelp()
}
