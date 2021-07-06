package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/testround/plans/example/pkg/dht"
)

// main for Standalone and debug run
func main() {
	fmt.Println("Start Bootstrap Host")
	// Shared cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, addr, err := dht.BootstrapDht(ctx) // empty peers for default bootstrapping
	// dht, err := dht.JoinDht(ctx, append(myPeers, ma))
	if err != nil {
		println("Could not start Bootstraper")
		panic(err)
	}

	println(id)
	println(addr.String())

	// save it to temp file to share it with others
	for {
		time.Sleep(time.Second)
	}
}
