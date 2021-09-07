package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/ihlec/bc_p2p/src/go/bc_dht/plans/coop_bc/pkg/dht"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"github.com/testground/sdk-go/runtime"
)

func UploadPeer(runenv *runtime.RunEnv, bootstrap_addr string) {
	// 1. Semantic Scholar Check

	apiURL := "https://api.semanticscholar.org/v1/paper/"

	// Get documentID by DOI

	sampleDocumentID := "77f59aac5011ae660181b6454a94c627d7339206"
	// cppd = 863f7197639325641f787caaf3a77a3f567fb24f
	// rbac = d7a3e44f86cb69dbc351b7d212312136ab6f0b8e
	// refs5 = 77f59aac5011ae660181b6454a94c627d7339206

	// Get all references by ID || What is the original work referencing?
	resp, err := http.Get(apiURL + sampleDocumentID)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var m dht.DocumentResponseStruct

	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}

	_ = m.References //Document References

	//inverted map || key:<referenceID> value:<[]paperIDsWhereItAppears>
	coCitationMap := make(map[string][]string, 1000) // Preallocate space for 1000 entries

	// Get citations by reference ID || Who else cited the references of the "original work"
	for _, refOfSubmission := range m.References {
		// Prepare the coCitationMap by setting submission's refernceIDs as key
		coCitationMap[refOfSubmission.PaperID] = nil

		// Get co-citations
		runenv.RecordMessage("Getting citations from submission-reference:")
		fmt.Printf("%+v\n", refOfSubmission)

		resp, err := http.Get(apiURL + refOfSubmission.PaperID)
		if err != nil {
			panic(err)
		}

		poteniallyCoCitedDocBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var cc dht.DocumentResponseStruct
		err = json.Unmarshal(poteniallyCoCitedDocBody, &cc)
		if err != nil {
			panic(err)
		}

		// Fill co-cite map
		// Saves all potenially-co-citing papersIDs to the inverted map with the cited reference's paperID as key
		coCitationMap[refOfSubmission.PaperID] = nil
		// Add all citing paperIDs
		for _, ccpID := range cc.References {
			coCitationMap[refOfSubmission.PaperID] = append(coCitationMap[refOfSubmission.PaperID], ccpID.PaperID)
		}

		runenv.RecordMessage("Co-Cit-Map-Entry: ")
		fmt.Printf("%+v\n", coCitationMap[refOfSubmission.PaperID])
	}
	runenv.RecordMessage("Entire Co-Cit-Map:")
	fmt.Printf("%+v\n", coCitationMap)

	// Create all k-combinations of the submission's references
	referencesDocIDs := make([]string, 0, len(coCitationMap))
	for key := range coCitationMap {
		referencesDocIDs = append(referencesDocIDs, key)
	}
	// We sort the strings to void redundant combinations || NO h(r1,r2) h(r2,r1) ONLY h(r1,r2)
	sort.Strings(referencesDocIDs)
	// TODO: this should not be hard-coded for k2 || allow for higher k in the future

	// Find elements for each ID combination
	combinations := [][]string{}
	for _, refIdA := range referencesDocIDs {
		for _, ref1dB := range referencesDocIDs {
			// Identical IDs would break combinational hash security
			if refIdA == ref1dB {
				continue
			}
			combinations = append(combinations, []string{refIdA, ref1dB})
		}
	}
	fmt.Printf("%+v\n", combinations)

	// Take each k-combination-element and create an intersection for the citing-documentIDs from the coCitationMap
	originalCombinations := [][]string{}
	for _, combination := range combinations {
		intersection := make([]string, 0)
		hash := make(map[interface{}]bool)

		for i := 0; i < len(coCitationMap[combination[0]]); i++ {
			el := coCitationMap[combination[0]][i]
			hash[el] = true
		}

		for i := 0; i < len(coCitationMap[combination[1]]); i++ {
			el := coCitationMap[combination[1]][i]
			if _, found := hash[el]; found {
				intersection = append(intersection, el)
			}
		}

		// Keep only unique and new combinations || check S2 for uniqe HDFs
		if len(intersection) <= 1 {
			originalCombinations = append(originalCombinations, combination)
		}
	}

	fmt.Printf("%+v\n", originalCombinations)

	// Filter public references
	fmt.Printf("%+v\n", len(combinations))
	fmt.Printf("%+v\n", len(originalCombinations))

	//fmt.Printf("%+v\n", m.References)
	//runenv.RecordMessage(m)

	runenv.RecordMessage("Join DHT")
	// New context for upload
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define bootstrap nodes
	ma, err := multiaddr.NewMultiaddr(bootstrap_addr)
	if err != nil {
		_ = ma
		panic(err)
	}

	var myPeers []multiaddr.Multiaddr
	dht, err := dht.JoinDht(ctx, runenv, append(myPeers, ma))
	if err != nil {
		runenv.RecordMessage("Could not join DHT")
		panic(err)
	}

	// Check the DHT for our hashed combinations
	inputDocHash := sha1.New()
	inputDocHash.Write([]byte(sampleDocumentID))

	var checkgroupInitial sync.WaitGroup
	channelForUnseenHashes := make(chan string, len(originalCombinations))
	for _, combination := range originalCombinations {
		h := sha1.New()
		h.Write([]byte(strings.Join(combination, "")))

		checkgroupInitial.Add(1)
		go func() {
			defer checkgroupInitial.Done()
			hashS := hex.EncodeToString(h.Sum(nil))
			if !check(ctx, runenv, dht, hashS) {
				channelForUnseenHashes <- hashS
			}
		}()
	}

	checkgroupInitial.Wait()

	// channel magic to identify last send message and read all
	messages := []string{}
	for {
		select {
		case msg := <-channelForUnseenHashes:
			messages = append(messages, msg)
			continue
		default:
		}
		break
	}
	close(channelForUnseenHashes)

	unseenHashes := []string{}
	for _, msg := range messages {
		unseenHashes = append(unseenHashes, msg)
		runenv.RecordMessage("Original New Hash: " + msg)
	}

	// 2. Batch UPLOAD in goroutine
	var uploadgroup sync.WaitGroup
	for _, element := range unseenHashes {
		uploadgroup.Add(1)
		go func(e string) {
			defer uploadgroup.Done()
			upload(ctx, runenv, dht, []string{sampleDocumentID, e})
		}(element)
	}
	uploadgroup.Wait()

	// // 3. Batch CHECK in goroutine || sanity check
	var checkgroup sync.WaitGroup
	for _, element := range unseenHashes {
		checkgroup.Add(1)
		go func(e string) {
			defer checkgroup.Done()
			check(ctx, runenv, dht, e)
		}(element)
	}
	checkgroup.Wait()

	// calc originality ratio RO
	runenv.RecordMessage("RO: " + fmt.Sprint(float32(len(unseenHashes))/float32(len(combinations))))

}

func upload(ctx context.Context, runenv *runtime.RunEnv, dht *kaddht.IpfsDHT, element []string) {
	fmt.Printf("PUT :: Document-Key: %s HDF: %s\n", element[0], element[1])
	err := dht.PutValue(ctx, "/v/"+element[1], []byte(element[0]))
	if err != nil {
		runenv.RecordMessage("Put Failed")
		panic(err)
	}
}

// check if HDF exists on DHT
func check(ctx context.Context, runenv *runtime.RunEnv, dht *kaddht.IpfsDHT, element string) bool {

	fmt.Printf("GET :: HDF: %q\n", element)
	myBytes, err := dht.GetValue(ctx, "/v/"+element)
	if err != nil {
		runenv.RecordMessage("GET Failed for: %q/n", element)
		return false

	} else {
		runenv.RecordMessage("Found HDF: " + element + " in DocumentID: " + string(myBytes[:]))
		return true
	}
}
