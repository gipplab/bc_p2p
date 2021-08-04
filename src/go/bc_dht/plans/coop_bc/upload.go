package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ihlec/bc_p2p/src/go/bc_dht/plans/coop_bc/pkg/dht"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/testground/sdk-go/runtime"
)

// main for Standalone and debug run
func UploadPeer(runenv *runtime.RunEnv, bootstrap_addr string) {
	runenv.RecordMessage("Join DHT")
	// New context for upload
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// b, err := ioutil.ReadFile("bootstrap_ID.tmp") // just pass the file name
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// bootstrap_addr := string(b)

	// Define bootstrap nodes
	// ma, err := multiaddr.NewMultiaddr("/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ") //ipfs
	ma, err := multiaddr.NewMultiaddr(bootstrap_addr)
	if err != nil {
		_ = ma
		panic(err)
	}

	var myPeers []multiaddr.Multiaddr
	// dht, err := dht.JoinDht(ctx, myPeers) // empty peers for default bootstrapping
	dht, err := dht.JoinDht(ctx, runenv, append(myPeers, ma))
	if err != nil {
		runenv.RecordMessage("Could not join DHT")
		panic(err)
	}

	// Single PUT GET to check network
	// dht.Provide() // TODO: might be more efficient
	// tryput:
	txValue := "valueDiesDAs"
	runenv.RecordMessage("PUT:", txValue)
	err = dht.PutValue(ctx, "/v/hello", []byte(txValue))
	if err != nil {
		runenv.RecordMessage("Put Failed")
		time.Sleep(time.Second)
		// goto tryput
	}
	myBytes, err := dht.GetValue(ctx, "/v/hello")
	rxValue := string(myBytes[:])
	runenv.RecordMessage("GET:", rxValue)
	if err != nil {
		runenv.RecordMessage("Get Failed")
		panic(err)
	} else {
		runenv.RecordMessage(rxValue)
	}

	// Batch UPLOAD in goroutine
	var uploadgroup sync.WaitGroup
	for _, element := range sampleData() {
		uploadgroup.Add(1)
		go upload(ctx, runenv, dht, element, &uploadgroup)
	}
	uploadgroup.Wait()

	// Batch CHECK in goroutine
	var checkgroup sync.WaitGroup
	for _, element := range sampleData() {
		checkgroup.Add(1)
		go check(ctx, runenv, dht, element, &checkgroup)
	}
	checkgroup.Wait()

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

func upload(ctx context.Context, runenv *runtime.RunEnv, dht *kaddht.IpfsDHT, element []string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("PUT :: Document-Key: %s HDF: %s\n", element[0], element[1])
	err := dht.PutValue(ctx, "/v/"+element[1], []byte(element[0]))
	if err != nil {
		runenv.RecordMessage("Put Failed")
		panic(err)
	}
}

func check(ctx context.Context, runenv *runtime.RunEnv, dht *kaddht.IpfsDHT, element []string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("GET :: HDF: %s\n", element[1])
	myBytes, err := dht.GetValue(ctx, "/v/"+element[1])
	if err != nil {
		runenv.RecordMessage("GET Failed")
		panic(err)
	} else {
		runenv.RecordMessage("Found HDF: " + element[1] + " in DocumentID: " + string(myBytes[:]))
	}
}
