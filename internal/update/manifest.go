package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuhaiin/kitte/manifest"
)

const (
	manifestRepository = "yuhaiin/kitte"
	manifestBranch     = "auto-update"
	manifestBaseURL    = "https://raw.githubusercontent.com/yuhaiin/kitte/auto-update/"
)

type manifestFileSpec struct {
	category   string
	kind       string
	usage      string
	listType   string
	format     string
	selectable bool
	sourceURL  string
}

func generateManifest(root string) error {
	files := make([]manifest.File, 0, 1024)

	add := func(rel string, spec manifestFileSpec) error {
		path := filepath.Join(root, filepath.FromSlash(rel))
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read manifest file %s: %w", rel, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat manifest file %s: %w", rel, err)
		}

		sum := sha256.Sum256(data)
		slashPath := filepath.ToSlash(rel)
		ext := filepath.Ext(slashPath)
		name := strings.TrimSuffix(filepath.Base(slashPath), ext)

		files = append(files, manifest.File{
			ID:         strings.TrimSuffix(slashPath, ext),
			Name:       name,
			Category:   spec.category,
			Kind:       spec.kind,
			Usage:      spec.usage,
			ListType:   spec.listType,
			Format:     spec.format,
			Path:       slashPath,
			URL:        manifestBaseURL + slashPath,
			SourceURL:  spec.sourceURL,
			Size:       info.Size(),
			SHA256:     hex.EncodeToString(sum[:]),
			Selectable: spec.selectable,
		})
		return nil
	}

	if err := add("geoip/Country.mmdb", manifestFileSpec{
		category: "geoip", kind: "database", usage: "maxminddb",
		format: "maxminddb", selectable: true, sourceURL: geoIPURL,
	}); err != nil {
		return err
	}
	if err := add("geosite/geosite.dat", manifestFileSpec{
		category: "geosite", kind: "database", usage: "source",
		format: "v2ray-geosite", selectable: false, sourceURL: geoSiteURL,
	}); err != nil {
		return err
	}

	if err := addDirectory(root, "geoip/geoip", nil, func(rel string) manifestFileSpec {
		return manifestFileSpec{
			category: "geoip", kind: "cidr", usage: "route-list",
			listType: "host", format: "text", selectable: true, sourceURL: geoIPURL,
		}
	}, add); err != nil {
		return err
	}

	if err := addDirectory(root, "geosite/geosite", nil, func(rel string) manifestFileSpec {
		return manifestFileSpec{
			category: "geosite", kind: "domain", usage: "route-list",
			listType: "host", format: "text", selectable: true, sourceURL: geoSiteURL,
		}
	}, add); err != nil {
		return err
	}

	sourceByOutput := make(map[string]string)
	for _, source := range ruleSources() {
		sourceByOutput[filepath.ToSlash(source.output)] = source.url
	}
	sourceByOutput["yuhaiin/tailscale.conf"] = tailscaleURL

	yuhaiinSkip := map[string]struct{}{
		"private.conf":    {},
		"yuhaiin_my.conf": {},
	}
	if err := addDirectory(root, "yuhaiin", yuhaiinSkip, func(rel string) manifestFileSpec {
		return manifestFileSpec{
			category: "yuhaiin", kind: "host", usage: "route-list",
			listType: "host", format: "text", selectable: true,
			sourceURL: sourceByOutput[filepath.ToSlash(rel)],
		}
	}, add); err != nil {
		return err
	}

	if err := addDirectory(root, "self", nil, func(rel string) manifestFileSpec {
		return manifestFileSpec{
			category: "self", kind: "host", usage: "route-list",
			listType: "host", format: "text", selectable: true,
		}
	}, add); err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].ID < files[j].ID })

	doc := manifest.Manifest{
		SchemaVersion: manifest.SchemaVersion,
		Repository:    manifestRepository,
		Branch:        manifestBranch,
		BaseURL:       manifestBaseURL,
		Files:         files,
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	data = append(data, '\n')

	if err := writeFileAtomic(filepath.Join(root, "manifest.json"), data, 0o644); err != nil {
		return fmt.Errorf("write manifest.json: %w", err)
	}
	return nil
}

func addDirectory(
	root, relDir string,
	skip map[string]struct{},
	spec func(rel string) manifestFileSpec,
	add func(rel string, spec manifestFileSpec) error,
) error {
	dir := filepath.Join(root, filepath.FromSlash(relDir))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read manifest directory %s: %w", relDir, err)
	}

	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		if _, ok := skip[entry.Name()]; ok {
			continue
		}
		rel := filepath.ToSlash(filepath.Join(relDir, entry.Name()))
		if err := add(rel, spec(rel)); err != nil {
			return err
		}
	}
	return nil
}
