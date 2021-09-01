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

	type DocumentResponseStruct struct {
		Abstract string `json:"abstract"`
		ArxivID  string `json:"arxivId"`
		Authors  []struct {
			AuthorID string `json:"authorId"`
			Name     string `json:"name"`
			URL      string `json:"url"`
		} `json:"authors"`
		CitationVelocity int `json:"citationVelocity"`
		Citations        []struct {
			ArxivID interface{} `json:"arxivId"`
			Authors []struct {
				AuthorID string `json:"authorId"`
				Name     string `json:"name"`
			} `json:"authors"`
			Doi           interface{}   `json:"doi"`
			Intent        []interface{} `json:"intent"`
			IsInfluential bool          `json:"isInfluential"`
			PaperID       string        `json:"paperId"`
			Title         string        `json:"title"`
			URL           string        `json:"url"`
			Venue         string        `json:"venue"`
			Year          int           `json:"year"`
		} `json:"citations"`
		CorpusID                 int      `json:"corpusId"`
		Doi                      string   `json:"doi"`
		FieldsOfStudy            []string `json:"fieldsOfStudy"`
		InfluentialCitationCount int      `json:"influentialCitationCount"`
		IsOpenAccess             bool     `json:"isOpenAccess"`
		IsPublisherLicensed      bool     `json:"isPublisherLicensed"`
		NumCitedBy               int      `json:"numCitedBy"`
		NumCiting                int      `json:"numCiting"`
		PaperID                  string   `json:"paperId"`
		References               []struct {
			ArxivID interface{} `json:"arxivId"`
			Authors []struct {
				AuthorID string `json:"authorId"`
				Name     string `json:"name"`
			} `json:"authors"`
			Doi           string   `json:"doi"`
			Intent        []string `json:"intent"`
			IsInfluential bool     `json:"isInfluential"`
			PaperID       string   `json:"paperId"`
			Title         string   `json:"title"`
			URL           string   `json:"url"`
			Venue         string   `json:"venue"`
			Year          int      `json:"year"`
		} `json:"references"`
		Title  string `json:"title"`
		Topics []struct {
			Topic   string `json:"topic"`
			TopicID string `json:"topicId"`
			URL     string `json:"url"`
		} `json:"topics"`
		URL   string `json:"url"`
		Venue string `json:"venue"`
		Year  int    `json:"year"`
	}
	var m DocumentResponseStruct

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

		var cc DocumentResponseStruct
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

	// TODO: Verify that keys and values are unique

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

		// Keep only unique and new combinations
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
	channelForUnseenHashes := make(chan []byte, sha1.New().Size())
	for _, combination := range originalCombinations {
		h := sha1.New()
		h.Write([]byte(strings.Join(combination, "")))

		checkgroupInitial.Add(1)
		go func() {
			if !check(ctx, runenv, dht, []string{string(inputDocHash.Sum(nil)), string(h.Sum(nil))}, &checkgroupInitial) {
				println("Unseen Hash: " + hex.EncodeToString([]byte(string(h.Sum(nil)))))

				channelForUnseenHashes <- h.Sum(nil)
			}
		}()
	}
	checkgroupInitial.Wait()

	messages := <-channelForUnseenHashes
	unseenHashes := make([]string, 0)
	for i := 0; i < len(messages); i += sha1.New().Size() {
		unseenHashes = append(unseenHashes, hex.EncodeToString(messages[i:i+sha1.New().Size()]))
		println("Als ob das nur einer ist!!: ", hex.EncodeToString(messages[i:i+sha1.New().Size()]))
	}

	// for _, uh := range *unseenHashes {
	// 	hex.EncodeToString(uh) //TODO: this should be a string, not a f* rune. The channel
	// }
	// runenv.RecordMessage(unseenHashes)
	runenv.RecordMessage(fmt.Sprint(len(unseenHashes)))

	// // 2. Batch UPLOAD in goroutine
	// var uploadgroup sync.WaitGroup
	// for _, element := range unseenHashes {
	// 	uploadgroup.Add(1)
	// 	go upload(ctx, runenv, dht, element, &uploadgroup)
	// }
	// uploadgroup.Wait()

	// // 3. Batch CHECK in goroutine
	// var checkgroup sync.WaitGroup
	// for _, element := range sampleData() {
	// 	checkgroup.Add(1)
	// 	go check(ctx, runenv, dht, element, &checkgroup)
	// }
	// checkgroup.Wait()

	// Sample DATA
	// // 2. Batch UPLOAD in goroutine
	// var uploadgroup sync.WaitGroup
	// for _, element := range sampleData() {
	// 	uploadgroup.Add(1)
	// 	go upload(ctx, runenv, dht, element, &uploadgroup)
	// }
	// uploadgroup.Wait()

	// // 3. Batch CHECK in goroutine
	// var checkgroup sync.WaitGroup
	// for _, element := range sampleData() {
	// 	checkgroup.Add(1)
	// 	go check(ctx, runenv, dht, element, &checkgroup)
	// }
	// checkgroup.Wait()
}

func sampleData() [][]string {
	// Todo: add fakeFile1000

	// Fake File
	fakeFile := [][]string{
		{"49001", "44b71eecde659689d848176615a2696aaeb2fb27"},
		{"49001", "5e84bace19194e76e9815b7e20df02d801089c7e"},
		{"49001", "f43b4c8ab78ab2cfaf6fbce63dac09087e352559"},
		{"49001", "67cc113d5a4c70e462c6e91db66eb90b996d4ec1"},
		{"49001", "8ab113d0d52cab6d190eb338dccbd07b39185420"},
		{"49001", "00e026e7bd70932a11711c80ce4ce235dda99860"},
		{"49001", "446c551127448c7fb3a4aa3799bada2195534fe6"},
		{"49001", "89975d4b6232f078fd3ae277da3aa75da20a0080"},
		{"49001", "cfb1a38a80d599420859793e40b3fc34cf537976"},
		{"49001", "3035020d7795f0fd041319b3208be61b09b8dcae"},
	}

	// // From CSV
	// csvfile, err := os.Open("test10_doc.csv")
	// if err != nil {
	// 	log.Fatalln("Couldn't open the csv file", err)
	// }
	// r := csv.NewReader(csvfile)
	// // Read each record from csv
	// record, err := r.ReadAll()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return fakeFile
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

func check(ctx context.Context, runenv *runtime.RunEnv, dht *kaddht.IpfsDHT, element []string, wg *sync.WaitGroup) bool {
	defer wg.Done()

	fmt.Printf("GET :: HDF: %q\n", hex.EncodeToString([]byte(element[1])))
	myBytes, err := dht.GetValue(ctx, "/v/"+element[1])
	if err != nil {
		runenv.RecordMessage("GET Failed for: %q/n", hex.EncodeToString([]byte(element[1])))
		return false

	} else {
		runenv.RecordMessage("Found HDF: " + hex.EncodeToString([]byte(element[1])) + " in DocumentID: " + string(myBytes[:]))
		return true
	}
}
