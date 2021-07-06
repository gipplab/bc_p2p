package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ipfs/testround/plans/example/pkg/dht"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

// main for Standalone and debug run
func main() {
	fmt.Println("Join DHT")
	// shared cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// define bootstrap nodes
	ma, err := multiaddr.NewMultiaddr("/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ")
	if err != nil {
		panic(err)
		_ = ma
	}
	var myPeers []multiaddr.Multiaddr
	//dht, err := dht.JoinDht(myPeers) // empty peers for default bootstrapping
	dht, err := dht.JoinDht(ctx, append(myPeers, ma))

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

	// Batch upload
	var uploadgroup sync.WaitGroup

	for i := 1; i <= 5; i++ {
		uploadgroup.Add(1)
		go uploader(ctx, dht, i, &uploadgroup)
	}

	uploadgroup.Wait()

}

func uploader(ctx context.Context, dht *kaddht.IpfsDHT, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Upload %d starting\n", id)
	// upload
	txValue := "valueDiesDAs"
	txKey := fmt.Sprint("/v/hello", id)
	println("PUT:", txValue)
	err := dht.PutValue(ctx, txKey, []byte(txValue))
	if err != nil {
		println("Put Failed")
		panic(err)
	}

	time.Sleep(time.Second)
	fmt.Printf("Upload %d done\n", id)
}
