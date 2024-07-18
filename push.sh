curl https://raw.githubusercontent.com/Loyalsoldier/geoip/release/Country.mmdb -o geoip/Country.mmdb
go generate ./geoip/...

# cd cn
# . ./get_cn.sh
# cd ..

cd yuhaiin
. ./get.sh
. ./update.sh
cd ..

git add .
git commit -m "update"
git push -u origin main
