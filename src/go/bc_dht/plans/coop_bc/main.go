// Welcome, testground plan writer!
// If you are seeing this for the first time, check out our documentation!
// https://app.gitbook.com/@protocol-labs/s/testground/

package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/testground/sdk-go/network"
	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
	"github.com/testground/sdk-go/sync"
)

//todo: followers all ready but bootstrapper is not reachable via DHT
//todo: ip4 address should not be localhost
//todo: check how to enforce a second LAN dht: https://github.com/libp2p/go-libp2p-kad-dht/blob/master/dual/dual_test.go
//todo: or use mdns to find default peers

func main() {
	run.Invoke(runf)
}

func runf(runenv *runtime.RunEnv) error {
	runenv.RecordMessage("Hello, Testground!")

	var (
		enrolledState = sync.State("enrolled")
		readyState    = sync.State("ready")
		releasedState = sync.State("released")

		addressesTopic = sync.NewTopic("bootstrapAddr", "")

		ctx = context.Background()
	)

	// instantiate a sync service client, binding it to the RunEnv.
	client := sync.MustBoundClient(ctx, runenv)
	defer client.Close()

	// instantiate a network client; see 'Traffic shaping' in the docs.
	netclient := network.NewClient(client, runenv)
	runenv.RecordMessage("waiting for network initialization")

	// wait for the network to initialize; this should be pretty fast.
	netclient.MustWaitNetworkInitialized(ctx)

	// signal entry in the 'enrolled' state, and obtain a sequence number.
	seq := client.MustSignalEntry(ctx, enrolledState)

	runenv.RecordMessage("my sequence ID: %d", seq)

	// if we're the first instance to signal, we'll become the LEADER.
	if seq == 1 {
		runenv.RecordMessage("i'm the bootstrapper.")
		numFollowers := runenv.TestInstanceCount - 1

		peerAddr, comProtocol, peerID := BootstrapPeer(ctx, runenv)

		// bootstrapper publishes its endpoint address on the 'addresses' topic
		seq := client.MustPublish(ctx, addressesTopic, peerAddr+comProtocol+peerID)

		runenv.RecordMessage("I am instance number %d publishing to the 'addresses' topic\n", seq)

		// ---------------------------------------
		// let's wait for the followers to signal.
		// ---------------------------------------
		runenv.RecordMessage("waiting for %d instances to become ready", numFollowers)
		err := <-client.MustBarrier(ctx, readyState, numFollowers).C
		if err != nil {
			return err
		}

		runenv.RecordMessage("the followers are all ready")
		runenv.RecordMessage("Lets upload...")

		UploadPeer(runenv, peerAddr+comProtocol+peerID)

		time.Sleep(1 * time.Second)
		runenv.RecordMessage("set...")
		time.Sleep(5 * time.Second)
		runenv.RecordMessage("go, release followers!")

		// signal on the 'released' state.
		client.MustSignalEntry(ctx, releasedState)

		time.Sleep(5 * time.Second)
		runenv.RecordMessage("Bootstrapper OUT...")

		client.Close()
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	sleep := rand.Intn(5) + 5
	runenv.RecordMessage("i'm a follower; signalling ready after %d seconds", sleep)
	time.Sleep(time.Duration(sleep) * time.Second)

	// consume all addresses from all peers
	ch := make(chan string)
	_ = client.MustSubscribe(ctx, addressesTopic, ch)

	addr := ""
	for i := 0; i < runenv.TestInstanceCount; i++ {
		addr = <-ch
		runenv.RecordMessage("received addr: %s", addr)
		if addr != "" {
			break
		}

	}
	IdlePeer(ctx, runenv, addr)

	runenv.RecordMessage("follower signalling now")

	// signal entry in the 'ready' state.
	client.MustSignalEntry(ctx, readyState)

	// wait until the leader releases us.
	err := <-client.MustBarrier(ctx, releasedState, 1).C
	if err != nil {
		return err
	}

	runenv.RecordMessage("i have been released")
	client.Close()
	return nil
}
