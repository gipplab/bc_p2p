package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht"
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

	return nil
}
