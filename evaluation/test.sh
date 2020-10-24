grep 'Duration' 1host2peers_batch_put_eth_run_10.txt | tail -n 1 | tr [:space:] '\n' | grep -v [a-z] | tr -d '{}\n' > temp.csv >> seconds_nanos.csv




# | sed -r '/^\s*$/d'
