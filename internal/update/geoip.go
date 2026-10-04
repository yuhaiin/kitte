package update

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oschwald/maxminddb-golang/v2"
)

type geoIPRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

func generateGeoIP(root string, database []byte) error {
	databasePath := filepath.Join(root, "geoip", "Country.mmdb")
	if err := writeFileAtomic(databasePath, database, 0o644); err != nil {
		return fmt.Errorf("write Country.mmdb: %w", err)
	}

	mdb, err := maxminddb.Open(databasePath)
	if err != nil {
		return fmt.Errorf("open Country.mmdb: %w", err)
	}
	defer mdb.Close()

	groups := make(map[string][]string)
	for network := range mdb.Networks() {
		var record geoIPRecord
		if err := network.Decode(&record); err != nil {
			return fmt.Errorf("decode %s: %w", network.Prefix(), err)
		}

		country := strings.ToUpper(strings.TrimSpace(record.Country.ISOCode))
		if country == "" {
			continue
		}
		groups[country] = append(groups[country], network.Prefix().String())
	}

	files := make(map[string][]byte, len(groups))
	for country, prefixes := range groups {
		sort.Strings(prefixes)
		var b bytes.Buffer
		for _, prefix := range prefixes {
			b.WriteString(prefix)
			b.WriteByte('\n')
		}
		files[country+".conf"] = b.Bytes()
	}

	if err := replaceDir(filepath.Join(root, "geoip", "geoip"), files); err != nil {
		return fmt.Errorf("write geoip rules: %w", err)
	}
	return nil
}
