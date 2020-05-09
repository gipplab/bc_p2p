#!/bin/bash

# Input number expected e.g. "./xstarttmux.sh 3"

# tmux kill-server
# tmux list-sessions
# tmux attach -t n

# Expects binary in same folder or in bin path

for (( i = 0; i <= $1; i++)) do tmux new-session -d -s $i && tmux send ./bc_p2p Enter && echo "$i Started"; done












