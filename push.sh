curl -L https://raw.githubusercontent.com/Loyalsoldier/geoip/release/Country.mmdb -o geoip/Country.mmdb
curl -L https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat -o geosite/geosite.dat

go generate ./geoip/...
go generate ./geosite/...

cd yuhaiin
. ./get.sh
. ./tailscale.sh
cd ..

# git add .
# git commit -m "update"
# git push -u origin main
