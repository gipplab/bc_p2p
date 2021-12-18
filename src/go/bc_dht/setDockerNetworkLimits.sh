
# sysctl -w net.core.somaxconn = 131072
# sysctl -w net.netfilter.nf_conntrack_max = 1048576
# sysctl -w net.ipv4.tcp_max_syn_backlog = 131072
# sysctl -w net.core.netdev_max_backlog = 524288
# sysctl -w net.ipv4.ip_local_port_range = 10000 65535
# sysctl -w net.ipv4.tcp_tw_recycle = 1
# sysctl -w net.ipv4.tcp_tw_reuse = 1
# sysctl -w net.core.rmem_max = 4194304
# sysctl -w net.core.wmem_max = 4194304
# sysctl -w net.ipv4.tcp_mem = 262144 524288 1572864
# sysctl -w net.ipv4.tcp_rmem = 16384 131072 4194304
# sysctl -w net.ipv4.tcp_wmem = 16384 131072 4194304
sysctl -w net.ipv4.neigh.default.gc_thresh2=4096
sysctl -w net.ipv4.neigh.default.gc_thresh3=32768