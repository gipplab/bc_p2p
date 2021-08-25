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

//todo: try all addresses of multiaddress for connection

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

	// if we're the first instance to signal, we'll become the BOOTSTRAPPER.
	if seq == 1 {
		runenv.RecordMessage("i'm the bootstrapper")

		peerAddr, comProtocol, peerID := BootstrapPeer(ctx, runenv)

		// publishes its endpoint address on the 'addresses' topic
		seq := client.MustPublish(ctx, addressesTopic, peerAddr+comProtocol+peerID)

		runenv.RecordMessage("I am instance number %d", seq)

		// signal entry in the 'ready' state for uploader
		client.MustSignalEntry(ctx, readyState)

		// wait until the uploader releases us.
		err := <-client.MustBarrier(ctx, releasedState, 1).C
		if err != nil {
			return err
		}

		runenv.RecordMessage("bootstrapper has been released")
		client.Close()

		return nil
	} else if seq == 2 {
		runenv.RecordMessage("i'm the uploader")

		// give bootstrapper some time
		time.Sleep(time.Duration(2) * time.Second)

		numFollowers := runenv.TestInstanceCount - 1

		runenv.RecordMessage("I am instance number %d \n", seq)

		// waiting for all other peers
		runenv.RecordMessage("waiting for %d instances to become ready", numFollowers)
		err := <-client.MustBarrier(ctx, readyState, numFollowers).C
		if err != nil {
			return err
		}

		runenv.RecordMessage("the followers are all ready")

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

		runenv.RecordMessage("Lets upload...")
		UploadPeer(runenv, addr)

		// signal on the 'released' state.
		runenv.RecordMessage("releasing peers!")
		client.MustSignalEntry(ctx, releasedState)

		runenv.RecordMessage("uploader has been closed")

		client.Close()
		return nil
	}

	runenv.RecordMessage("i'm a normal peer")
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	time.Sleep(time.Duration(n+5) * time.Second)

	// consume all addresses from all peers
	ch := make(chan string)
	_ = client.MustSubscribe(ctx, addressesTopic, ch)

	addr := ""
	for i := 0; i < runenv.TestInstanceCount; i++ {
		addr = <-ch
		runenv.RecordMessage("received addr: %s", addr)
		if addr != "" {
			IdlePeer(ctx, runenv, addr)
			break
		}

	}

	runenv.RecordMessage("follower signalling now")

	// signal entry in the 'ready' state.
	client.MustSignalEntry(ctx, readyState)

	// wait until the uploader releases us.
	err := <-client.MustBarrier(ctx, releasedState, 1).C
	if err != nil {
		return err
	}

	runenv.RecordMessage("i have been released")
	client.Close()
	return nil
}
