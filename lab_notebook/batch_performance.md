# Batch Performance Measures

Content protecting bibliographic coupling requires the average processing of 30 references per document. 
A content protecting combinations of hashes by k leads to 30^k references.
For k = 2 this leads to around 900 hashes per document.
Additionally, in future work we plan to pilot to approach with at least 3 hosting institutions. 

We therefore measure the PUT and GET performance of the DHT with 1000 hashes on: 
1. 2 to 10 peers on a single host  
2. 3 peers and using 3 hosts

## 1. Single Host Test (17.05.2020)
We use a simple shell script to create and kill peers using tmux (terminal multiplexer). 

The interacting peer is run manually and logged using 
`$ script test_n.log`
to store the output logs.

Starting additional peers:
`$ sh xstarttmux.sh n`

Killing peers:
`$ tmux kill-server`
or 
`$ tmux kill-session -t n`

List peers:
`$ tmux list-sessions`

### 2 Peers
`$ sh xstarttmux.sh 0`
`$ script test_2.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 20:00:57.786470 +02:00
Finish: 2020-05-17 20:01:00.641616 +02:00

GET: 
Start:  2020-05-17 20:01:28.851195 +02:00
Finish: 2020-05-17 20:01:28.873104 +02:00

It is not possible to query without at least one other peer running

### 3 Peers
`$ sh xstarttmux.sh 1`
`$ script test_3.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 20:10:13.913990 +02:00
Finish: 2020-05-17 20:10:18.766391 +02:00

GET form running peer: 
Start:  2020-05-17 20:10:32.582171 +02:00
Finish: 2020-05-17 20:10:32.603894 +02:00

GET form restarted peer: 
Start:  2020-05-17 20:11:05.903699 +02:00
Finish: 2020-05-17 20:11:06.401548 +02:00

We are restarting the interacting peer to retrieve true network queries from here on

### 4 Peers
`$ sh xstarttmux.sh 2`
`$ script test_4.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 20:22:07.923557 +02:00
Finish: 2020-05-17 20:22:13.746870 +02:00

An immediately GET call after restart causes:
`Failed to get record: NotFound { key: Key(b"640f26174733738a543c9f195f211da321492f39"), closest_peers: [] }`
This is probably caused by an empty routing table, which is filled after a few seconds.

GET form restarted peer: 
Start:  2020-05-17 20:23:58.487065 +02:00
Finish: 2020-05-17 20:23:59.214879 +02:00

### 5 Peers
`$ sh xstarttmux.sh 3`
`$ script test_5.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 20:30:11.798906 +02:00
Finish: 2020-05-17 20:30:16.721987 +02:00

An immediately GET call after restart causes:
Failed to get record: NotFound { key: Key(b"640f26174733738a543c9f195f211da321492f39"), closest_peers: [] }
This is probably caused by an empty routing table, which is filled after a few seconds.

GET form restarted peer: 
Start:  2020-05-17 20:31:01.557491 +02:00
Finish: 2020-05-17 20:31:02.331466 +02:00

### 6 Peers
`$ sh xstarttmux.sh 4`
`$ script test_6.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 20:35:45.329804 +02:00
Finish: 2020-05-17 20:35:48.814688 +02:00

GET form restarted peer: 
Start:  2020-05-17 20:36:29.827893 +02:00
Finish: 2020-05-17 20:36:30.669686 +02:00

### 7 Peers
`$ sh xstarttmux.sh 5`
`$ script test_7.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 20:40:18.516880 +02:00
Finish: 2020-05-17 20:40:21.905937 +02:00

GET form restarted peer: 
Start:  2020-05-17 20:41:18.372146 +02:00
Finish: 2020-05-17 20:41:19.202370 +02:00

### 8 Peers
`$ sh xstarttmux.sh 6`
`$ script test_8.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 20:56:23.872652 +02:00
Finish: 2020-05-17 20:56:30.714599 +02:00

GET form restarted peer: 
Start:  2020-05-17 20:57:08.150285 +02:00
Finish: 2020-05-17 20:57:08.972635 +02:00

### 9 Peers
`$ sh xstarttmux.sh 7`
`$ script test_9.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 21:00:20.194054 +02:00
Finish: 2020-05-17 21:00:33.948413 +02:00

GET form restarted peer: 
Start:  2020-05-17 21:01:29.124415 +02:00
Finish: 2020-05-17 21:01:29.938341 +02:00

### 10 Peers
`$ sh xstarttmux.sh 8`
`$ script test_10.log`
`$ ./bc_p2p`
`$ BATCH ../../data/1000_k1.csv`
`$ CHECK ../../data/1000_k1.csv`
`$ ./bc_p2p` (kill and re-run)
`$ CHECK ../../data/1000_k1.csv`
`$ exit `

PUT
Start:  2020-05-17 21:04:25.473117 +02:00
Finish: 2020-05-17 21:04:36.366399 +02:00

GET form restarted peer: 
Start:  2020-05-17 21:05:10.811768 +02:00
Finish: 2020-05-17 21:05:11.620212 +02:00