curl --silent https://controlplane.tailscale.com/derpmap/default | jq -r '.Regions[].Nodes[] | "\(.HostName)\n\(.IPv4)\n\(.IPv6)"' > tailscale.conf
echo "login.tailscale.com" >> tailscale.conf
echo "controlplane.tailscale.com" >> tailscale.conf
echo "log.tailscale.com" >> tailscale.conf
echo "192.200.0.0/24" >> tailscale.conf
echo "2606:B740:49::/48" >> tailscale.conf

# https://tailscale.com/kb/1082/firewall-ports
# https://tailscale.com/kb/1232/derp-servers