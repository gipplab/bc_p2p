package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	yamux "github.com/libp2p/go-libp2p-yamux"
	tcp "github.com/libp2p/go-tcp-transport"
	"github.com/testground/sdk-go/runtime"
)

// Demonstrate test output functions
// This method emits two Messages and one Metric
// func ExampleOutput(runenv *runtime.RunEnv) error {
// 	runenv.RecordMessage("Hello, World.")
// 	runenv.RecordMessage("Additional arguments: %d", len(runenv.TestInstanceParams))
// 	runenv.R().RecordPoint("donkeypower", 3.0)
// 	return nil
// }

func DhtIdle(runenv *runtime.RunEnv) error {
	runenv.RecordMessage("Idle...")
	runenv.RecordMessage("Additional arguments: %d", len(runenv.TestInstanceParams))

	// init libp2p connector
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
	)

	// allow parallel connections on the same transport
	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	// set transport to listen on ipv4
	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
	)
	if err != nil {
		panic(err)
	}

	for _, addr := range host.Addrs() {
		fmt.Println("Listening on", addr)
	}

	//DHT
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	// searching peers
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// manually look for bootstrap peers
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
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
