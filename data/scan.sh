#!/bin/bash

# input integers 
# $1 as number of hosts
# $2 as number of peers
# $3 as number of runs 
tmux new-session -s query10 -d
tmux detach
sleep 1s
# initial BATCH upload
tmux send-keys -t 0 "./bc_p2p" ENTER
sleep 3s
tmux send-keys -t 0 "BATCH /Users/corihle/GIT/bc_p2p/data/1000_k1.csv" ENTER
sleep 60s
tmux send-keys -t 0 C-c

# scan runs
for (( i = 0; i <= $3; i++)) do
tmux send-keys -t 00 "script $1host$2peers_scan_run$i.txt" ENTER
sleep 1s
tmux send-keys -t 0 "./bc_p2p" ENTER
sleep 3s
tmux send-keys -t 0 "CHECK /Users/corihle/GIT/bc_p2p/data/1000_k1.csv" ENTER
sleep 10s
tmux send-keys -t 0 C-c
sleep 1s
tmux send-keys -t 0 "exit" ENTER
sleep 1s;
done
tmux kill-server
