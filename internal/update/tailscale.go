package update

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type derpMap struct {
	Regions map[string]struct {
		Nodes []struct {
			HostName string `json:"HostName"`
			IPv4     string `json:"IPv4"`
			IPv6     string `json:"IPv6"`
		} `json:"Nodes"`
	} `json:"Regions"`
}

func generateTailscale(root string, data []byte) error {
	var m derpMap
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("decode DERP map: %w", err)
	}

	keys := make([]string, 0, len(m.Regions))
	for key := range m.Regions {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, aErr := strconv.Atoi(keys[i])
		b, bErr := strconv.Atoi(keys[j])
		if aErr == nil && bErr == nil {
			return a < b
		}
		return keys[i] < keys[j]
	})

	lines := make([]string, 0, len(keys)*6)
	seen := map[string]struct{}{}
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		lines = append(lines, value)
	}

	for _, key := range keys {
		for _, node := range m.Regions[key].Nodes {
			add(node.HostName)
			add(node.IPv4)
			add(node.IPv6)
		}
	}

	for _, value := range []string{
		"login.tailscale.com",
		"controlplane.tailscale.com",
		"log.tailscale.com",
		"192.200.0.0/24",
		"2606:B740:49::/48",
	} {
		add(value)
	}

	if err := writeFileAtomic(filepath.Join(root, "yuhaiin", "tailscale.conf"), joinLines(lines), 0o644); err != nil {
		return fmt.Errorf("write tailscale.conf: %w", err)
	}
	return nil
}
