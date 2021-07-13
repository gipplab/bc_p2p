package main

import (
	"fmt"

	"github.com/testground/sdk-go/run"
	"github.com/testground/sdk-go/runtime"
)

func main() {
	run.Invoke(runf)
}

// Pick a different example function to run
// depending on the name of the test case.
func runf(runenv *runtime.RunEnv) {
	switch c := runenv.TestCase; c {
	case "idle":
		IdlePeer() // Will never quit
	// case "upload":
	// 	return DhtBatchUpload(runenv)
	case "bootstarp":
		uploadPeer()
	case "panic":
		ExamplePanic(runenv)
	case "params":
		ExampleParams(runenv)
	case "sync":
		ExampleSync(runenv)
	case "metrics":
		ExampleMetrics(runenv)
	case "artifact":
		ExampleArtifact(runenv)
	default:
		msg := fmt.Sprintf("Unknown Testcase %s", c)
		println(msg)
	}
}
