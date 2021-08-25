sysctl net.bridge.bridge-nf-call-iptables=0
sysctl net.bridge.bridge-nf-call-arptables=0
sysctl net.bridge.bridge-nf-call-ip6tables=0
sysctl net.ipv4.conf.all.forwarding=1
sudo iptables -P FORWARD ACCEPT
