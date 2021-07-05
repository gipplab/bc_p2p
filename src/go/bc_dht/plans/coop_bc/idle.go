package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	tcp "github.com/libp2p/go-tcp-transport"
)

// main for Standalone and debug run
// func main() {
// 	fmt.Println("Idle Mode")
// 	DhtIdle()
// }

//// TestPlan Specific Logging
// Demonstrate test output functions
// This method emits two Messages and one Metric
// func ExampleOutput(runenv *runtime.RunEnv) error {
// 	runenv.RecordMessage("Hello, World.")
// 	runenv.RecordMessage("Additional arguments: %d", len(runenv.TestInstanceParams))
// 	runenv.R().RecordPoint("donkeypower", 3.0)
// 	return nil
// }

func DhtIdle() error {
	// JOIN DHT
	// shared cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set ipv4 listener
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
	)
	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)
	// muxers := libp2p.ChainOptions(
	// 	libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
	// 	libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	// ) // is default

	// create peer (libp2p node)
	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		//muxers,
	)
	if err != nil {
		panic(err)
	}

	// Show Local IPv4
	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}

	// Init the DHT
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	// Clean routing table
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Look for other peers
	var wg sync.WaitGroup // create goroutines group
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		// go through IPFS default peers
		// TODO: compare public IPFS with Private DHT
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()

			// Ping
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Println(err)
			} else {
				log.Println("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return nil
}
