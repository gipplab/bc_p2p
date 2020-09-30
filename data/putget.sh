#!/bin/bash

# input integers 
# $1 as number of hosts
# $2 as number of peers
# $3 as number of runs 
tmux new-session -s putget -d
sleep 1s

# Log
tmux send-keys -t 0 "script $1host$2peers_putget_run$i.txt" ENTER

# Runs
for (( i = 1; i <= $3; i++)) do

# PUT
tmux send-keys -t 0 "./bc_p2p" ENTER
sleep 3s
tmux send-keys -t 0 "PUT test value" ENTER
sleep 3s
tmux send-keys -t 0 C-c

# GET
sleep 1s
tmux send-keys -t 0 "./bc_p2p" ENTER
sleep 3s
tmux send-keys -t 0 "GET test" ENTER
sleep 3s
tmux send-keys -t 0 C-c
sleep 1s
done
tmux send-keys -t 0 "exit" ENTER
sleep 1s;
tmux kill-server
