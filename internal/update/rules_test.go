package update

import (
	"strings"
	"testing"
)

func TestTransformDNSMasq(t *testing.T) {
	input := []byte("server=/example.cn/114.114.114.114\r\n# comment\r\n")

	got, err := transformDNSMasq(true)(input)
	if err != nil {
		t.Fatal(err)
	}
	if want := "*.example.cn\n# comment\n"; string(got) != want {
		t.Fatalf("wildcard transform = %q, want %q", got, want)
	}

	got, err = transformDNSMasq(false)(input)
	if err != nil {
		t.Fatal(err)
	}
	if want := "example.cn\n# comment\n"; string(got) != want {
		t.Fatalf("plain transform = %q, want %q", got, want)
	}
}

func TestTransformAdblockListDoesNotTreatHeaderAsRegexClass(t *testing.T) {
	input := []byte("[Adblock Plus 2.0]\n*.adblock.example\nclock.example\n||tracker.example^\n")
	got, err := transformAdblockList(input)
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	for _, want := range []string{"*.adblock.example", "clock.example", "||tracker.example^"} {
		if !strings.Contains(text, want+"\n") {
			t.Fatalf("result %q is missing %q", text, want)
		}
	}
	if strings.Contains(text, "[Adblock") {
		t.Fatalf("header was not removed: %q", text)
	}
}

func TestTransformHosts(t *testing.T) {
	input := []byte("# header\n0.0.0.0 ads.example # comment\n127.0.0.1 tracker.example tracker2.example\nplain.example\n::1 localhost\n")
	got, err := transformHosts(input)
	if err != nil {
		t.Fatal(err)
	}
	if want := "ads.example\ntracker.example\ntracker2.example\nplain.example\n"; string(got) != want {
		t.Fatalf("transformHosts = %q, want %q", got, want)
	}
}

func TestTransformAntiAD(t *testing.T) {
	input := []byte("# comment\nexample.com # inline\n*.already.example\nexample.com\n")
	got, err := transformAntiAD(input)
	if err != nil {
		t.Fatal(err)
	}
	if want := "*.example.com\n*.already.example\n"; string(got) != want {
		t.Fatalf("transformAntiAD = %q, want %q", got, want)
	}
}
