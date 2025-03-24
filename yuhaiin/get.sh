APPLE_CHINA_DOMAIN=https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/apple.china.conf
GOOGLE_CHINA_DOMAIN=https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/google.china.conf
CHINA_DOMAIN=https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf
AD_WARS_HOSTS=https://raw.githubusercontent.com/jdlingyu/ad-wars/master/hosts
PGLYOYO="https://pgl.yoyo.org/adservers/serverlist.php?hostformat=adblock&showintro=0&mimetype=plaintext"
ANTI_AD_DOMAINS="https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt"
damengzhu_banad="https://raw.githubusercontent.com/damengzhu/banad/refs/heads/main/hosts.txt"
VRChat_Analytics_Blocker="https://raw.githubusercontent.com/DubyaDude/VRChat-Analytics-Blocker/refs/heads/master/blocklist/hosts.txt"

curl ${CHINA_DOMAIN} | sed $'s/\r$//' | sed 's/server=\//\*\./g' | sed 's/\/114\.114\.114\.114//g' > accelerated-domains.china.conf
curl ${GOOGLE_CHINA_DOMAIN} | sed $'s/\r$//' | sed 's/server=\///g' | sed 's/\/114\.114\.114\.114//g' > google.china.conf
curl ${APPLE_CHINA_DOMAIN} | sed $'s/\r$//' | sed 's/server=\///g' | sed 's/\/114\.114\.114\.114//g' > apple.china.conf


curl "${ANTI_AD_DOMAINS}" | sed $'s/\r$//' | sed '/^#/d' | sed 's/#.*$//g' | sed 's/# /*\./g' | sed '/^ *$/d' | sed 's/^/*./g' > anti-ad-domains.txt
curl "${PGLYOYO}" | sed $'s/\r$//' | sed '/ *[Adblock] */d' > pglyoyo.txt
curl ${AD_WARS_HOSTS} | sed $'s/\r$//' | sed 's/127.0.0.1 //g' | sed 's/0.0.0.0 //g' | sed '/^#.*$/'d | sed '1,2d' | sed '/^ *$/d' > ad_wars_hosts
curl ${VRChat_Analytics_Blocker} | sed $'s/\r$//' | sed 's/127.0.0.1 //g' | sed 's/0.0.0.0 //g' | sed '/^#.*$/'d | sed '1,2d' | sed '/^ *$/d' > VRChat_Analytics_Blocker
curl ${damengzhu_banad} | sed $'s/\r$//' | sed 's/127.0.0.1 //g' | sed 's/0.0.0.0 //g' | sed '/^#.*$/'d | sed '1,2d' | sed '/^ *$/d' > damengzhu_banad
