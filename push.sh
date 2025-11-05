curl https://raw.githubusercontent.com/Loyalsoldier/geoip/release/Country.mmdb -o geoip/Country.mmdb
go generate ./geoip/...

cd yuhaiin
. ./get.sh
. ./tailscale.sh
cd ..

# git add .
# git commit -m "update"
# git push -u origin main
