package update

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

func generateGeoSite(root string, data []byte) error {
	if err := writeFileAtomic(filepath.Join(root, "geosite", "geosite.dat"), data, 0o644); err != nil {
		return fmt.Errorf("write geosite.dat: %w", err)
	}

	var siteList routercommon.GeoSiteList
	if err := proto.Unmarshal(data, &siteList); err != nil {
		return fmt.Errorf("decode geosite.dat: %w", err)
	}

	groups := make(map[string][]string)
	seen := make(map[string]map[string]struct{})

	for _, site := range siteList.Entry {
		name := strings.ToLower(strings.TrimSpace(site.CountryCode))
		if name == "" {
			continue
		}
		if seen[name] == nil {
			seen[name] = make(map[string]struct{})
		}

		for _, domain := range site.Domain {
			value := strings.TrimSpace(domain.Value)
			if value == "" || domain.Type == routercommon.Domain_Regex {
				continue
			}
			if domain.Type == routercommon.Domain_RootDomain {
				value = "*." + value
			}
			if _, ok := seen[name][value]; ok {
				continue
			}
			seen[name][value] = struct{}{}
			groups[name] = append(groups[name], value)
		}
	}

	files := make(map[string][]byte, len(groups))
	filenames := make([]string, 0, len(groups))
	for name, domains := range groups {
		filename := name + ".conf"
		filenames = append(filenames, filename)

		var b bytes.Buffer
		for _, domain := range domains {
			b.WriteString(domain)
			b.WriteByte('\n')
		}
		files[filename] = b.Bytes()
	}

	if err := replaceDir(filepath.Join(root, "geosite", "geosite"), files); err != nil {
		return fmt.Errorf("write geosite rules: %w", err)
	}

	sort.Strings(filenames)
	var readme bytes.Buffer
	readme.WriteString("# Geosite\n\n")
	readme.WriteString("Generated from `geosite.dat` by `go run ./cmd/kitte update`.\n\n")
	readme.WriteString("| File | Raw link |\n| --- | --- |\n")
	const baseURL = "https://raw.githubusercontent.com/yuhaiin/kitte/auto-update/geosite/geosite"
	for _, filename := range filenames {
		fmt.Fprintf(&readme, "| `%s` | [%s](%s/%s) |\n", filename, filename, baseURL, filename)
	}

	if err := writeFileAtomic(filepath.Join(root, "geosite", "README.md"), readme.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write geosite README: %w", err)
	}
	return nil
}
