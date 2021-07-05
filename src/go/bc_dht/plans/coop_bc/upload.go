package main

import (
	"context"
	"fmt"

	"github.com/ipfs/testround/plans/example/pkg/dht"
	"github.com/multiformats/go-multiaddr"
)

// main for Standalone and debug run
func main() {
	fmt.Println("Join DHT")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start peer
	ma, err := multiaddr.NewMultiaddr("/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ")
	if err != nil {
		panic(err)
		_ = ma
	}
	var myPeers []multiaddr.Multiaddr
	dht, err := dht.JoinDht(myPeers) //dht.JoinDht(append(myPeers, ma))

	// upload
	// dht.Provide() // TODO: might be more efficient
	txValue := "valueDiesDAs"
	println("PUT:", txValue)
	err = dht.PutValue(ctx, "/v/hello", []byte(txValue))
	if err != nil {
		println("Put Failed")
		panic(err)
	}

	myBytes, err := dht.GetValue(ctx, "/v/hello")
	rxValue := string(myBytes[:])
	println("GET:", rxValue)
	if err != nil {
		println("Get Failed")
		panic(err)
	} else {
		println(rxValue)
	}

}

// func DhtBatchUpload(runenv *runtime.RunEnv) error {
// 	runenv.RecordMessage("Uploading...")

// 	runenv.RecordMessage("Additional arguments: %d", len(runenv.TestInstanceParams))
// 	return nil
// }
