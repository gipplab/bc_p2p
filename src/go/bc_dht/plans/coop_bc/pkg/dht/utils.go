package dht

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	tcp "github.com/libp2p/go-tcp-transport"
	"github.com/multiformats/go-multiaddr"
)

// JoinDht start the node and tries to connect to the provided bootstrapPeers.
// If no bootstrapPeers are probided, the default IPFS bootstrapPeers are used.
func JoinDht(bootstrapPeers []multiaddr.Multiaddr) error {
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

	// create peer (libp2p node)
	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
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
	if len(bootstrapPeers) == 0 {
		bootstrapPeers = dht.DefaultBootstrapPeers
	}

	for _, peerAddr := range bootstrapPeers {
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
