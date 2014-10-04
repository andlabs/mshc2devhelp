// 4 october 2014
package main

type Book struct {
	XMLNS		string		`xml:"xmlns,attr"`
	Title			string		`xml:"title,attr"`
	Link			string		`xml:"link,attr"`
	Author		string		`xml:"author,attr"`
	Name		string		`xml:"name,attr"`
	Version		string		`xml:"version,attr"`
	Language		string		`xml:"language,attr"`
	Chapters		Chapters		`xml:"chapters"`
}

type Chapters struct {
	Sub			[]*Sub		`xml:"sub"`
}

type Sub struct {
	Name		string		`xml:"name,attr"`
	Link			string		`xml:"link,attr"`
	Sub			[]*Sub		`xml:"sub"`
}

var book = &Book{
	XMLNS:		"http://www.devhelp.net/book",
	Title:			"Windows Desktop App Development",		// TODO
	Version:		"2",
	Language:	"c",
}

func toSub(entry *Entry) *Sub {
	s := &Sub{
		Name:	entry.Name,
	}
	for _, e := range entry.Children {
		s.Sub = append(s.Sub, toSub(e))
	}
	return s
}

func buildDevhelp(bookname string) {
	book.Name = bookname
	for _, e := range toplevels {
		book.Chapters.Sub = append(book.Chapters.Sub, toSub(e))
	}
}
