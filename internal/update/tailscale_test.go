package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateTailscaleSortsRegionsAndSkipsEmptyValues(t *testing.T) {
	root := t.TempDir()
	data := []byte(`{"Regions":{"10":{"Nodes":[{"HostName":"derp10","IPv4":"10.0.0.10","IPv6":""}]},"2":{"Nodes":[{"HostName":"derp2","IPv4":"10.0.0.2","IPv6":"::2"}]}}}`)

	if err := generateTailscale(root, data); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(root, "yuhaiin", "tailscale.conf"))
	if err != nil {
		t.Fatal(err)
	}
	wantPrefix := "derp2\n10.0.0.2\n::2\nderp10\n10.0.0.10\n"
	if len(got) < len(wantPrefix) || string(got[:len(wantPrefix)]) != wantPrefix {
		t.Fatalf("tailscale output = %q, want prefix %q", got, wantPrefix)
	}
}
