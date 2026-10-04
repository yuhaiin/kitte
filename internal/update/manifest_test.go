package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yuhaiin/kitte/manifest"
)

func TestGenerateManifest(t *testing.T) {
	root := t.TempDir()

	write := func(rel, value string) {
		t.Helper()
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(value), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("geoip/Country.mmdb", "mmdb")
	write("geoip/geoip/JP.conf", "1.0.16.0/20\n")
	write("geosite/geosite.dat", "dat")
	write("geosite/geosite/google.conf", "*.google.com\n")
	write("yuhaiin/tailscale.conf", "login.tailscale.com\n")
	write("yuhaiin/private.conf", "must-not-leak.example\n")
	write("self/block.conf", "blocked.example\n")

	if err := generateManifest(root); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}

	var got manifest.Manifest
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.SchemaVersion != manifest.SchemaVersion {
		t.Fatalf("schema version = %d, want %d", got.SchemaVersion, manifest.SchemaVersion)
	}
	if got.BaseURL != manifestBaseURL {
		t.Fatalf("base URL = %q", got.BaseURL)
	}

	if len(got.Files) != 6 {
		t.Fatalf("file count = %d, want 6: %+v", len(got.Files), got.Files)
	}
	for i := 1; i < len(got.Files); i++ {
		if got.Files[i-1].ID >= got.Files[i].ID {
			t.Fatalf("manifest is not sorted: %q >= %q", got.Files[i-1].ID, got.Files[i].ID)
		}
	}

	var jp *manifest.File
	for i := range got.Files {
		if got.Files[i].Path == "yuhaiin/private.conf" {
			t.Fatal("ignored private file leaked into manifest")
		}
		if got.Files[i].Path == "geoip/geoip/JP.conf" {
			jp = &got.Files[i]
		}
	}
	if jp == nil {
		t.Fatal("JP geoip file missing")
	}
	if jp.ListType != "host" || !jp.Selectable || jp.Usage != "route-list" {
		t.Fatalf("unexpected JP metadata: %+v", *jp)
	}
	sum := sha256.Sum256([]byte("1.0.16.0/20\n"))
	if jp.SHA256 != hex.EncodeToString(sum[:]) || jp.Size != int64(len("1.0.16.0/20\n")) {
		t.Fatalf("unexpected JP integrity metadata: %+v", *jp)
	}
}
