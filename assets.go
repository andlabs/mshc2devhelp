// 4 october 2014
package main

import (
	"io"
	"io/ioutil"
	"bytes"
)

type Asset struct {
	MSHC	string
	Data		[]byte
}

var assets = make(map[string]*Asset)

func addAsset(mshcname string, name string, r io.Reader) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)		// TODO
	}
	a := &Asset{
		MSHC:	mshcname,
		Data:	b,
	}
	if assets[name] == nil {
		assets[name] = a
		return
	}
	if !bytes.Equal(assets[name].Data, a.Data) {
		panic("duplicate differing assets " + name + ": " + assets[name].MSHC + " vs " + a.MSHC)
	}
}
