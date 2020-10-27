#!/bin/bash

# input integers 
# $1 as number of hosts
# $2 as number of peers
# $3 as number of runs 
# $4 as name of configuration (e.g. wifi)

# Monitoring
# tmux attach -t putget
tmux new-session -s batchPUT -d
sleep 1s

# Runs
for (( i = 1; i <= $3; i++)) do
    # Log start
    tmux send-keys -t 0 "script ./results/$1host$2peers_batch_put_$4_run_$i.txt" ENTER

    # PUT
    tmux send-keys -t 0 "./bc_p2p" ENTER
    sleep 5s
    tmux send-keys -t 0 "BATCH 1000_k1.csv" ENTER
    sleep 60s
    tmux send-keys -t 0 C-c

    # Log end
    tmux send-keys -t 0 "exit" ENTER
    sleep 1s;

    # Search for last duration an save to .csv
    grep 'Duration' ./results/$1host$2peers_batch_put_$4_run_$i.txt | tail -n 1 | tr [:space:] '\n' | grep -v [a-z] | tr -d '{}\n' >> ./results/timings$1host$2peers_batch_put_$4.csv

    echo '' >> ./results/timings$1host$2peers_batch_put_$4.csv;
done

tmux kill-server
