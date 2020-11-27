package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	bucketName = "discoverrudy.appspot.com"
	appspotURL = "https://firebasestorage.googleapis.com/v0/b/discoverrudy.appspot.com/o/static"
)

var (
	regionID string
	verbose  bool
	test     bool
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
	flag.BoolVar(&test, "test", true, "whether to upload to test collections in Firestore")

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

func main() {
	flag.Parse()

	if regionID == "" {
		log.Fatalln("compress: error: regionID is empty")
	}

	zipFilePath := "compressed/" + regionID + ".zip"
	fileInfo, err := os.Stat(zipFilePath)
	if os.IsNotExist(err) {
		log.Fatalf("upload: datafile archive %s doesn't exist\n", zipFilePath)
	}

	regionName := "TODO"

	featured, err := parseFeatured(regionID)
	if err != nil {
		log.Fatalln("upload: error parsing featured:", err)
	}

	storagePrefix := regionID
	datafilesCollection := "datafiles"
	if test {
		storagePrefix += "Test"
		datafilesCollection += "Test"
	}

	fileURL := fmt.Sprintf("%s/%s/%s?alt=media", appspotURL, storagePrefix, fileInfo.Name())
	thumbURL := fmt.Sprintf("%s/%s/thumb.webp?alt=media", appspotURL, storagePrefix)
	thumbMiniURL := fmt.Sprintf("%s/%s/thumb_mini.webp?alt=media", appspotURL, storagePrefix)

	thumbBlurhash, err := makeThumbBlurhash(regionID)
	if err != nil {
		log.Fatalln("upload: error making a blurhash:", err)
	}

	fmt.Printf(thumbBlurhash)
	os.Exit(69)

	meta := Datafile{
		Available:        true,
		Featured:         featured,
		FileSize:         fileInfo.Size(),
		FileURL:          url.QueryEscape(fileURL),
		LastUploadedTime: time.Time{},
		Position:         1, // TODO: Handle position
		RegionID:         regionID,
		RegionName:       regionName,
		IsTestVersion:    test,
		ThumbBlurhash:    thumbBlurhash,
		ThumbMiniURL:     url.QueryEscape(thumbMiniURL),
		ThumbURL:         url.QueryEscape(thumbURL),
	}

	_, err = firestoreClient.Collection(datafilesCollection).Doc(regionID).Set(ctx, meta)
	if err != nil {
		log.Fatalf("error updating document %#v in /datafiles: %v\n", regionID, err)
	}
}
