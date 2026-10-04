# kitte

`kitte` builds the rule files consumed by Yuhaiin from several upstream data sources.
The generated files are maintained on the [`auto-update`](https://github.com/yuhaiin/kitte/tree/auto-update) branch by GitHub Actions.

## Layout

- `cmd/kitte/` — updater command.
- `internal/update/` — downloads, parsers, generators, and tests.
- `geoip/` — MaxMind-compatible country database plus generated per-country CIDR files.
- `geosite/` — V2Ray geosite database plus generated per-category domain files.
- `yuhaiin/` — generated domain/IP lists used by Yuhaiin.
- `self/` — manually maintained local rules; the updater does not overwrite them.

The generated output paths intentionally stay stable so existing raw URLs and downstream users do not need to change.

## Update locally

Requirements: Go version from `go.mod`. No `curl`, `sed`, `jq`, or shell-specific behavior is required.

```sh
go test ./...
go run ./cmd/kitte update
```

Use `-root` when running against another checkout:

```sh
go run ./cmd/kitte -root /path/to/kitte update
```

The updater validates HTTP status codes, retries transient download failures, writes files through temporary files, and makes generated output deterministic where upstream ordering is not meaningful.

## Sources

| Output | Upstream |
| --- | --- |
| `geoip/` | [Loyalsoldier/geoip](https://github.com/Loyalsoldier/geoip) |
| `geosite/` | [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) |
| China domain lists | [felixonmars/dnsmasq-china-list](https://github.com/felixonmars/dnsmasq-china-list) |
| `anti-ad-domains.txt` | [privacy-protection-tools/anti-AD](https://github.com/privacy-protection-tools/anti-AD) |
| `ad_wars_hosts` | [jdlingyu/ad-wars](https://github.com/jdlingyu/ad-wars) |
| `damengzhu_banad` | [damengzhu/banad](https://github.com/damengzhu/banad) |
| `VRChat_Analytics_Blocker` | [DubyaDude/VRChat-Analytics-Blocker](https://github.com/DubyaDude/VRChat-Analytics-Blocker) |
| `pglyoyo.txt` | [Peter Lowe's blocklist](https://pgl.yoyo.org/adservers/) |
| `tailscale.conf` | [Tailscale DERP map](https://controlplane.tailscale.com/derpmap/default) |

## Automation

`.github/workflows/update.yml` runs the Go updater daily and keeps the long-lived `auto-update` pull request/branch refreshed. The branch is useful as the latest generated dataset while `main` keeps the updater source and manually maintained rules reviewable.
