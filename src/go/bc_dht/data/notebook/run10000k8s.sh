#!/bin/bash

echo "If more than 30 peers are used, make sure to fix the network config:" 
echo "sysctl -w net.ipv4.neigh.default.gc_thresh2=4096"
echo "sysctl -w net.ipv4.neigh.default.gc_thresh3=32768" 
echo ""
echo "Make sure the docker image is build before the batch run for accurate results!"
echo ""

# $1 max number of peers
# $2 increase step per run
# $3 number of repetitions
if [ $# -eq 0 ]
  then
    echo "No arguments supplied. 1: max number of peers; 2: increase step per run; 3: number of repetitions"
    exit 1
fi

pushd ../outputs/aws_k8s

# Repetitions
for (( r = 1; r <= $3; r++)) do
    echo "Repetition $r of $3: "

    # Runs
    for (( i = 0; i <= $1; i = i+$2)) do
        if [ $i == 0 ]; then 
            continue
        fi
        echo "Run with $i peers. Results at: "
        pwd
        message=$(testground run single --plan=coopbc --testcase=bc --runner=cluster:k8s --builder=docker:go --instances=$i --collect --build-cfg go_proxy_mode=direct)
        runid=$(echo $message | grep -o '....................$')
        echo $runid
        sleep 120s
        testground collect --runner=cluster:k8s $runid
        sleep 5s
        tar -xzvf $runid.tgz
        rm $runid.tgz
    done
    sleep 60s
done
popd 
