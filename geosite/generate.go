//go:build ignore

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Read the raw bytes from the file
	data, err := os.ReadFile("geosite.dat")
	if err != nil {
		fmt.Printf("Failed to read geosite.dat: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal the protobuf data
	var siteList routercommon.GeoSiteList
	if err := proto.Unmarshal(data, &siteList); err != nil {
		fmt.Printf("Failed to unmarshal geosite.dat: %v\n", err)
		os.Exit(1)
	}

	// Clean up existing output directory
	outputDir := "geosite"
	os.RemoveAll(outputDir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Iterate over each GeoSite entry
	for _, site := range siteList.Entry {
		countryCode := strings.ToLower(site.CountryCode)
		if countryCode == "" {
			continue
		}

		filename := fmt.Sprintf("%s/%s.conf", outputDir, countryCode)
		f, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Failed to create file %s: %v\n", filename, err)
			continue
		}

		for _, domain := range site.Domain {
			// We only write the domain value.
			// Types like Regex, Plain, Domain, Full are collapsed into just the string.
			if domain.Value != "" {
				fmt.Fprintln(f, domain.Value)
			}
		}
		f.Close()
	}
}
