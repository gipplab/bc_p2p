package main

import (
	"context"

	"github.com/ihlec/bc_p2p/src/go/bc_dht/plans/coop_bc/pkg/dbc"
	"github.com/testground/sdk-go/runtime"
)

// main for Standalone and debug run
// returns host.addr, protocol, peerID
func BootstrapPeer(ctx context.Context, runenv *runtime.RunEnv) (string, string, string) {
	runenv.RecordMessage("Start Bootstrap Host")
	// Shared cancelable context
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	dht, err := dbc.BootstrapDht(ctx, runenv) // empty peers for default bootstrapping
	// dht, err := dht.JoinDht(ctx, append(myPeers, ma))
	if err != nil {
		runenv.RecordMessage("Could not start Bootstraper")
		panic(err)
	}

	runenv.RecordMessage(dht.PeerID().String())
	runenv.RecordMessage(dht.Host().Addrs()[0].String()) // 0 Localhost, 1 LAN, 2 WAN

	// // write PeerID to file
	// f, err := os.Create("bootstrap_ID.tmp")
	// _, err = f.WriteString(dht.Host().Addrs()[0].String() + "/p2p/" + dht.PeerID().String()) //"/ip4/192.168.2.70/tcp/45009/p2p/QmasJcoadBm2LQ1WTtxJbnaywnHhRYiboTWSKDZCBGytTm"
	// f.Close()
	// if err != nil {
	// 	runenv.RecordMessage("Could not write File")
	// 	panic(err)
	// }

	// // save it to temp file to share it with others
	// for {
	// 	time.Sleep(time.Second)
	// }

	return dht.Host().Addrs()[0].String(), "/p2p/", dht.PeerID().String() // 0 Localhost, 1 LAN, 2 WAN
}
