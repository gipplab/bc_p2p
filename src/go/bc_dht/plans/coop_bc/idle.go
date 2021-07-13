package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ipfs/testround/plans/example/pkg/dht"
	"github.com/multiformats/go-multiaddr"
)

// Start and join a peer in idle mode
func IdlePeer() {
	fmt.Println("Join DHT")
	// Shared cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := ioutil.ReadFile("bootstrap_ID.tmp") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	bootstrap_addr := string(b)

	// Define bootstrap nodes
	//ma, err := multiaddr.NewMultiaddr("/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ") //ipfs
	ma, err := multiaddr.NewMultiaddr(bootstrap_addr) //custom
	if err != nil {
		_ = ma
		panic(err)
	}

	var myPeers []multiaddr.Multiaddr
	//dht, err := dht.JoinDht(ctx, myPeers) // empty peers for default bootstrapping
	dht, err := dht.JoinDht(ctx, append(myPeers, ma))
	if err != nil {
		println("Could not join DHT")
		panic(err)
	}

	println("Own PeerID: " + dht.PeerID().String())
	for {
		time.Sleep(time.Second)
	}
}
