package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/testround/plans/example/pkg/dht"
	"github.com/multiformats/go-multiaddr"
)

// main for Standalone and debug run
func main() {
	fmt.Println("Join DHT")
	// Shared cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define bootstrap nodes
	//ma, err := multiaddr.NewMultiaddr("/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ") //ipfs
	ma, err := multiaddr.NewMultiaddr("/ip4/192.168.2.70/tcp/45009/p2p/QmasJcoadBm2LQ1WTtxJbnaywnHhRYiboTWSKDZCBGytTm") //custom
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
