//go:build ignore

package main

import (
	"fmt"
	"os"
	"sort"
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

	var filenames []string

	// Iterate over each GeoSite entry
	for _, site := range siteList.Entry {
		countryCode := strings.ToLower(site.CountryCode)
		if countryCode == "" {
			continue
		}

		filename := fmt.Sprintf("%s.conf", countryCode)
		filePath := fmt.Sprintf("%s/%s", outputDir, filename)
		f, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Failed to create file %s: %v\n", filePath, err)
			continue
		}

		filenames = append(filenames, filename)

		for _, domain := range site.Domain {
			// We only write the domain value.
			// Types like Regex, Plain, Domain, Full are collapsed into just the string.
			if domain.Value != "" {
				fmt.Fprintln(f, domain.Value)
			}
		}
		f.Close()
	}

	generateReadme(outputDir, filenames)
}

func generateReadme(outputDir string, filenames []string) {
	sort.Strings(filenames)

	f, err := os.Create("README.md")
	if err != nil {
		fmt.Printf("Failed to create README.md: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "# Geosite")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "| File | Raw Link |")
	fmt.Fprintln(f, "| --- | --- |")

	baseURL := "https://raw.githubusercontent.com/yuhaiin/kitte/auto-update/geosite/geosite"

	for _, filename := range filenames {
		fmt.Fprintf(f, "| %s | [%s](%s/%s) |\n", filename, filename, baseURL, filename)
	}
}
