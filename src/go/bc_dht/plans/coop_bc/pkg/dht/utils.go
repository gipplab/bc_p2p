package dht

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	record "github.com/libp2p/go-libp2p-record"
	tcp "github.com/libp2p/go-tcp-transport"
	"github.com/multiformats/go-multiaddr"
)

type blankValidator struct{}

func (blankValidator) Validate(_ string, _ []byte) error        { return nil }
func (blankValidator) Select(_ string, _ [][]byte) (int, error) { return 0, nil }

// JoinDht start the node and tries to connect to the provided bootstrapPeers.
// If no bootstrapPeers are probided, the default IPFS bootstrapPeers are used.
func JoinDht(ctx context.Context, bootstrapPeers []multiaddr.Multiaddr) (*kaddht.IpfsDHT, error) {

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
	kademliaDHT, err := kaddht.New(ctx, host, kaddht.Mode(kaddht.ModeServer))
	// !!! ModeServer will cause troubles when running outside a private network behind NAT
	if err != nil {
		panic(err)
	}

	// // Clean routing table
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	kademliaDHT.Validator.(record.NamespacedValidator)["v"] = blankValidator{} // Value (v)

	// Look for other peers
	var wg sync.WaitGroup // create goroutines group
	if len(bootstrapPeers) == 0 {
		bootstrapPeers = kaddht.DefaultBootstrapPeers
	}

	for _, peerAddr := range bootstrapPeers {
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

	return kademliaDHT, nil
}

func BootstrapDht(ctx context.Context) (*kaddht.IpfsDHT, error) {
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

	fmt.Println("Bootstrap PeerID: " + host.ID().String())

	// Init the DHT
	kademliaDHT, err := kaddht.New(ctx, host, kaddht.Mode(kaddht.ModeServer))
	if err != nil {
		panic(err)
	}

	kademliaDHT.Validator.(record.NamespacedValidator)["v"] = blankValidator{} // might not be needed

	return kademliaDHT, nil
}
