package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ipfs/testround/plans/example/pkg/dht"
)

// main for Standalone and debug run
func bootstrapPeer() {
	fmt.Println("Start Bootstrap Host")
	// Shared cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dht, err := dht.BootstrapDht(ctx) // empty peers for default bootstrapping
	// dht, err := dht.JoinDht(ctx, append(myPeers, ma))
	if err != nil {
		println("Could not start Bootstraper")
		panic(err)
	}

	println(dht.PeerID().String())
	println(dht.Host().Addrs()[0].String())

	// write PeerID to file
	f, err := os.Create("bootstrap_ID.tmp")
	_, err = f.WriteString(dht.Host().Addrs()[0].String() + "/p2p/" + dht.PeerID().String()) //"/ip4/192.168.2.70/tcp/45009/p2p/QmasJcoadBm2LQ1WTtxJbnaywnHhRYiboTWSKDZCBGytTm"
	f.Close()
	if err != nil {
		println("Could not write File")
		panic(err)
	}

	// save it to temp file to share it with others
	for {
		time.Sleep(time.Second)
	}
}
