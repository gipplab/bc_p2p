package main

import (
	"context"

	"github.com/ihlec/bc_p2p/src/go/bc_dht/plans/coop_bc/pkg/dbc"
	"github.com/multiformats/go-multiaddr"
	"github.com/testground/sdk-go/runtime"
)

// Start and join a peer in idle mode
func IdlePeer(ctx context.Context, runenv *runtime.RunEnv, bootstrap_addr string) {
	runenv.RecordMessage("Join DHT")
	// Shared cancelable context
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// b, err := ioutil.ReadFile("bootstrap_ID.tmp") // just pass the file name
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// bootstrap_addr := string(b)

	// Define bootstrap nodes
	//ma, err := multiaddr.NewMultiaddr("/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ") //ipfs
	ma, err := multiaddr.NewMultiaddr(bootstrap_addr) //custom
	if err != nil {
		_ = ma
		panic(err)
	}

	var myPeers []multiaddr.Multiaddr
	//dht, err := dbc.JoinDht(ctx, myPeers) // empty peers for default bootstrapping
	dht, err := dbc.JoinDht(ctx, runenv, append(myPeers, ma))
	if err != nil {
		runenv.RecordMessage("Could not join DHT")
		panic(err)
	}

	runenv.RecordMessage("Own PeerID: " + dht.PeerID().String())
	// for {
	// 	time.Sleep(time.Second)
	// }
}
