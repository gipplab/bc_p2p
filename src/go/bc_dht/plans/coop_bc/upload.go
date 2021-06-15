package main

import (
	"github.com/testground/sdk-go/runtime"
)

func DhtBatchUpload(runenv *runtime.RunEnv) error {
	runenv.RecordMessage("Uploading...")

	runenv.RecordMessage("Additional arguments: %d", len(runenv.TestInstanceParams))
	return nil
}
