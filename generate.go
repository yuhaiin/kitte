//go:build ignore

package main

import (
	"bytes"
	"compress/gzip"
	"os"
)

//go:generate go run generate.go

func main() {
	data, err := os.ReadFile("yuhaiin/yuhaiin.conf")
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer(nil)

	gw := gzip.NewWriter(b)
	gw.Write(data)
	gw.Close()

	err = os.WriteFile("bypass.gz", b.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
