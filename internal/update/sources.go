package update

const (
	geoIPURL     = "https://raw.githubusercontent.com/Loyalsoldier/geoip/release/Country.mmdb"
	geoSiteURL   = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
	tailscaleURL = "https://controlplane.tailscale.com/derpmap/default"
)

type textSource struct {
	name      string
	url       string
	output    string
	transform func([]byte) ([]byte, error)
}

func ruleSources() []textSource {
	return []textSource{
		{
			name:      "dnsmasq accelerated domains",
			url:       "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf",
			output:    "yuhaiin/accelerated-domains.china.conf",
			transform: transformDNSMasq(true),
		},
		{
			name:      "dnsmasq Google China domains",
			url:       "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf",
			output:    "yuhaiin/google.china.conf",
			transform: transformDNSMasq(false),
		},
		{
			name:      "dnsmasq Apple China domains",
			url:       "https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf",
			output:    "yuhaiin/apple.china.conf",
			transform: transformDNSMasq(false),
		},
		{
			name:      "anti-AD domains",
			url:       "https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt",
			output:    "yuhaiin/anti-ad-domains.txt",
			transform: transformAntiAD,
		},
		{
			name:      "PGLYoyo ad servers",
			url:       "https://pgl.yoyo.org/adservers/serverlist.php?hostformat=adblock&showintro=0&mimetype=plaintext",
			output:    "yuhaiin/pglyoyo.txt",
			transform: transformAdblockList,
		},
		{
			name:      "ad-wars hosts",
			url:       "https://raw.githubusercontent.com/jdlingyu/ad-wars/master/hosts",
			output:    "yuhaiin/ad_wars_hosts",
			transform: transformHosts,
		},
		{
			name:      "VRChat analytics blocker",
			url:       "https://raw.githubusercontent.com/DubyaDude/VRChat-Analytics-Blocker/refs/heads/master/blocklist/hosts.txt",
			output:    "yuhaiin/VRChat_Analytics_Blocker",
			transform: transformHosts,
		},
		{
			name:      "damengzhu banad",
			url:       "https://raw.githubusercontent.com/damengzhu/banad/refs/heads/main/hosts.txt",
			output:    "yuhaiin/damengzhu_banad",
			transform: transformHosts,
		},
	}
}
