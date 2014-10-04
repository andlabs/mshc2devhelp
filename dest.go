// 4 october 2014
package main

import (
	"fmt"
	"os"
	"io"
	"archive/zip"
	"path/filepath"
)

func copyOne(mshcname string, in string, out string) {
	z, err := zip.OpenReader(mshcname)
	if err != nil {
		panic(err)
	}
	defer z.Close()
	for _, f := range z.File {
		if f.Name != in {
			continue
		}
		r, err := f.Open()
		if err != nil {
			panic(err)		// TODO
		}
		of, err := os.Create(out)
		if err != nil {
			panic(err)		// TODO
		}
		_, err = io.Copy(of, r)
		if err != nil {
			panic(err)		// TODO
		}
		of.Close()
		r.Close()
		return
	}
	panic(in + " not found in " + mshcname)
}

func buildDestinationFolder(dir string) {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		panic(err)		// TODO
	}
	i := 0
	for _, e := range entries {
		e.Dest = filepath.Join(dir, fmt.Sprintf("%d.htm", i))
		copyOne(e.MSHC, e.File, e.Dest)
		i++
	}
	copyAssets(dir)
}
