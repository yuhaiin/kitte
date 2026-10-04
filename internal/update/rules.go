package update

import (
	"bytes"
	"fmt"
	"net/netip"
	"path/filepath"
	"strings"
)

func updateRules(root string, downloads map[string][]byte) error {
	for _, source := range ruleSources() {
		raw, ok := downloads[source.url]
		if !ok {
			return fmt.Errorf("missing downloaded source %s", source.name)
		}
		data, err := source.transform(raw)
		if err != nil {
			return fmt.Errorf("transform %s: %w", source.name, err)
		}
		if err := writeFileAtomic(filepath.Join(root, source.output), data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", source.output, err)
		}
	}
	return nil
}

func transformDNSMasq(wildcard bool) func([]byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		lines := splitLines(data)
		out := make([]string, 0, len(lines))
		for _, line := range lines {
			if strings.HasPrefix(line, "server=/") {
				rest := strings.TrimPrefix(line, "server=/")
				if i := strings.LastIndexByte(rest, '/'); i >= 0 {
					rest = rest[:i]
				}
				if wildcard {
					rest = "*." + rest
				}
				line = rest
			}
			out = append(out, line)
		}
		return joinLines(out), nil
	}
}

func transformAntiAD(data []byte) ([]byte, error) {
	lines := splitLines(data)
	out := make([]string, 0, len(lines))
	seen := make(map[string]struct{}, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "*.") {
			line = "*." + line
		}
		if _, ok := seen[line]; ok {
			continue
		}
		seen[line] = struct{}{}
		out = append(out, line)
	}
	return joinLines(out), nil
}

func transformAdblockList(data []byte) ([]byte, error) {
	lines := splitLines(data)
	out := make([]string, 0, len(lines))
	seen := make(map[string]struct{}, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(strings.ToLower(line), "[adblock") {
			continue
		}
		if _, ok := seen[line]; ok {
			continue
		}
		seen[line] = struct{}{}
		out = append(out, line)
	}
	return joinLines(out), nil
}

func transformHosts(data []byte) ([]byte, error) {
	lines := splitLines(data)
	out := make([]string, 0, len(lines))
	seen := make(map[string]struct{}, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		start := 0
		if _, err := netip.ParseAddr(fields[0]); err == nil {
			start = 1
		}
		for _, host := range fields[start:] {
			host = strings.TrimSpace(strings.TrimSuffix(host, "."))
			if host == "" || host == "localhost" || host == "localhost.localdomain" {
				continue
			}
			if _, ok := seen[host]; ok {
				continue
			}
			seen[host] = struct{}{}
			out = append(out, host)
		}
	}
	return joinLines(out), nil
}

func splitLines(data []byte) []string {
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	data = bytes.ReplaceAll(data, []byte("\r"), []byte("\n"))
	text := strings.TrimSuffix(string(data), "\n")
	if text == "" {
		return nil
	}
	return strings.Split(text, "\n")
}

func joinLines(lines []string) []byte {
	if len(lines) == 0 {
		return nil
	}
	return []byte(strings.Join(lines, "\n") + "\n")
}
