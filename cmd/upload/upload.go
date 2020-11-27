package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	bucketName = "discoverrudy.appspot.com"
)

var (
	regionID string
	verbose  bool
)

var (
	ctx             = context.Background()
	firestoreClient *firestore.Client
	storageClient   *storage.Client
)

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.BoolVar(&verbose, "verbose", false, "true for extensive logging")
	flag.Parse()

	if regionID == "" {
		log.Fatalln("compress: error: regionID is empty")
	}

	opt := option.WithCredentialsFile("./key.json")

	var err error
	firestoreClient, err = firestore.NewClient(ctx, "discoverrudy", opt)
	if err != nil {
		log.Fatalf("upload: error initializing firestore: %v\n", err)
	}

	storageClient, err = storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalf("upload: error initializing storage: %v\n", err)
	}
}

// parseMeta parses metadata for the datafile from the database.
func parseMeta() (*Meta, error) {
	contribsFile, err := os.Open("database/" + regionID + "/meta/contributors.json")
	if err != nil {
		// log.Fatalln("upload: error opening contributors file:", err)
		return nil, err
	}
	defer contribsFile.Close()

	b, err := ioutil.ReadAll(contribsFile)
	if err != nil {
		// log.Fatalln("upload: error reading from contributors file:", err)
		return nil, err
	}

	var contributors []string
	err = json.Unmarshal(b, &contributors)
	if err != nil {
		// log.Fatalln("upload: error unmarshalling contributors file:", err)
		return nil, err
	}

	fmt.Printf("contributors: %v\n", contributors)

	return nil, errors.New("not implemented yet")
}

func main() {
	zipFilePath := "compressed/" + regionID + ".zip"
	fileInfo, err := os.Stat(zipFilePath)
	if os.IsNotExist(err) {
		log.Fatalf("upload: datafile archive %s doesn't exist\n", zipFilePath)
	}

	fileInfo.Size()

	meta, err := parseMeta()
	if err != nil {
		log.Fatalln("upload: error parsing metadata:", err)
	}

	_ = meta.Sources
}
