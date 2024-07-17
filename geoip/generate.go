//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang"
)

func main() {
	mdb, err := maxminddb.Open("Country.mmdb")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.RemoveAll("geoip")

	record := struct {
		Country struct {
			IsoCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}{}

	nets := mdb.Networks(maxminddb.SkipAliasedNetworks)

	for nets.Next() {
		subnet, err := nets.Network(&record)
		if err != nil {
			continue
		}

		name := strings.ToUpper(record.Country.IsoCode)

		f := getFiles(name)

		fmt.Fprintln(f, subnet)
	}

	for _, f := range fm {
		f.Close()
	}
}

var fm = map[string]*os.File{}

func getFiles(isCode string) *os.File {
	f, ok := fm[isCode]
	if ok {
		return f
	}

	err := os.MkdirAll("geoip", 0755)
	if err != nil {
		panic(err)
	}

	f, err = os.Create("geoip/" + isCode + ".conf")
	if err != nil {
		panic(err)
	}

	fm[isCode] = f
	return f
}
