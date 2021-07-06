package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ipfs/testround/plans/example/pkg/dht"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

// main for Standalone and debug run
func main() {
	fmt.Println("Join DHT")
	// Shared cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define bootstrap nodes
	ma, err := multiaddr.NewMultiaddr("/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ")
	if err != nil {
		panic(err)
		_ = ma
	}
	var myPeers []multiaddr.Multiaddr
	dht, err := dht.JoinDht(ctx, myPeers) // empty peers for default bootstrapping
	// dht, err := dht.JoinDht(ctx, append(myPeers, ma))

	// Single PUT GET to check network
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

	// Batch upload in goroutine
	var uploadgroup sync.WaitGroup

	// Upload
	for _, element := range sampleData() {
		uploadgroup.Add(1)
		go uploader(ctx, dht, element, &uploadgroup)
	}

	uploadgroup.Wait()

}

func sampleData() [][]string {
	// From CSV
	csvfile, err := os.Open("test10_doc.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)
	// Read each record from csv
	record, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return record
}

func uploader(ctx context.Context, dht *kaddht.IpfsDHT, element []string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("PUT :: Document-Key: %s HDF: %s\n", element[0], element[1])
	err := dht.PutValue(ctx, "/v/"+element[0], []byte(element[1]))
	if err != nil {
		println("Put Failed")
		panic(err)
	}
}
