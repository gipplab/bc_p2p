#!/bin/bash

# input integers 
# $1 as number of hosts
# $2 as number of peers
# $3 as number of runs 
# $4 as name of configuration (e.g. wifi)

# Monitoring
# tmux attach -t putget

tmux new-session -s batchGET -d
sleep 1s

# PUT
tmux send-keys -t 0 "./bc_p2p" ENTER
sleep 3s
tmux send-keys -t 0 "BATCH /home/workbook/GIT/bc_p2p/data/1000_k1.csv" ENTER
sleep 60s
tmux send-keys -t 0 C-cc

# Runs
for (( i = 1; i <= $3; i++)) do

# Log start
tmux send-keys -t 0 "script $1host$2peers_batch_get_$4_run_$i.txt" ENTER

# GET
tmux send-keys -t 0 "./bc_p2p" ENTER
sleep 3s
tmux send-keys -t 0 "CHECK /home/workbook/GIT/bc_p2p/data/1000_k1.csv" ENTER
sleep 15s
tmux send-keys -t 0 C-c

# Log end
tmux send-keys -t 0 "exit" ENTER
sleep 1s

# Search for last duration an save to .csv
grep 'Duration' $1host$2peers_batch_get_$4_run_$i.txt | tail -n 1 | tr [:space:] '\n' | grep -v [a-z] | tr -d '{}\n' | sed -r '/^\s*$/d' >> seconds_nanos.csv
echo '' >> timings$1host$2peers_batch_get_$4seconds_nanos.csv
done

tmux kill-server


