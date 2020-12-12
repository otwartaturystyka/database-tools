package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/bartekpacia/database-tools/internal"
	"github.com/pkg/errors"

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
	lang     string
	position int
	onlyMeta bool
	verbose  bool
	noTest   bool
)

var (
	firestoreClient *firestore.Client
	storageClient   *storage.Client
)

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.StringVar(&lang, "lang", "pl", "language of the datafile to upload")
	flag.IntVar(&position, "position", 1, "position at which the datafile should show in the app")
	flag.BoolVar(&onlyMeta, "only-meta", false, "true to upload only metadata (not the .zip file)")
	flag.BoolVar(&noTest, "no-test", false, "true to upload to *production* collection in Firestore")
	flag.BoolVar(&verbose, "verbose", false, "true for extensive logging")

	opt := option.WithCredentialsFile("./key.json")

	var err error
	firestoreClient, err = firestore.NewClient(context.Background(), "discoverrudy", opt)
	if err != nil {
		log.Fatalf("upload: error initializing firestore: %v\n", err)
	}

	storageClient, err = storage.NewClient(context.Background(), opt)
	if err != nil {
		log.Fatalf("upload: error initializing storage: %v\n", err)
	}
}

func main() {
	flag.Parse()

	if regionID == "" {
		log.Fatalln("compress: regionID is empty")
	}

	zipFilePath := "compressed/" + regionID + ".zip"
	fileInfo, err := os.Stat(zipFilePath)
	if os.IsNotExist(err) {
		log.Fatalf("upload: datafile archive %s doesn't exist\n", zipFilePath)
	}

	regionName := "TODO"

	meta, err := parseMeta(regionID, lang)
	if err != nil {
		log.Fatalln("upload: error parsing meta:", err)
	}

	storagePrefix := regionID
	datafilesCollection := "datafiles"
	if !noTest {
		storagePrefix += "Test"
		datafilesCollection += "Test"
	}

	fileURL := fmt.Sprintf("%s/%s/%s?alt=media", appspotURL, storagePrefix, fileInfo.Name())
	thumbURL := fmt.Sprintf("%s/%s/thumb.webp?alt=media", appspotURL, storagePrefix)
	thumbMiniURL := fmt.Sprintf("%s/%s/thumb_mini.webp?alt=media", appspotURL, storagePrefix)

	fmt.Printf("upload: begin making thumb blurhash...")
	thumbBlurhash, err := makeThumbBlurhash(regionID)
	if err != nil {
		log.Fatalln("upload: error making a blurhash:", err)
	}
	fmt.Println("ok")

	fmt.Println("upload: you are going to upload a data pack with the following metadata")

	datafileData := internal.FirestoreDatafile{
		Available:        true,
		Featured:         meta.Featured,
		FileSize:         fileInfo.Size(),
		FileURL:          fileURL,
		LastUploadedTime: time.Time{},
		Position:         1, // TODO: Handle position
		RegionID:         regionID,
		RegionName:       regionName,
		IsTestVersion:    !noTest,
		ThumbBlurhash:    thumbBlurhash,
		ThumbMiniURL:     thumbMiniURL,
		ThumbURL:         thumbURL,
	}

	datafileDataJSON, err := json.MarshalIndent(datafileData, "", "  ")
	if err != nil {
		log.Fatalln("upload: failed to marshal datafileData to JSON:", err)
	}
	fmt.Println(string(datafileDataJSON))

	fmt.Println("upload: continue? (Y/n)")
	accepted, err := askForConfirmation()
	if err != nil {
		log.Fatalf("upload: failed to get response: %v\n", err)
	}

	if !accepted {
		fmt.Println("upload: operation canceled by the user")
		os.Exit(0)
	}

	// Upload compressed datafile
	if !onlyMeta {
		func() {
			localPath := "compressed/" + regionID + ".zip"
			cloudPath := "static/" + storagePrefix + "/rudy.zip"
			upload(localPath, cloudPath, "application/zip")
		}()
	}

	// Upload thumb
	func() {
		localPath := "database/" + regionID + "/meta/thumb.webp"
		cloudPath := "static/" + storagePrefix + "/thumb.webp"
		upload(localPath, cloudPath, "image/webp")
	}()

	// Upload minified thumb
	func() {
		localPath := "database/" + regionID + "/meta/thumb_mini.webp"
		cloudPath := "static/" + storagePrefix + "/thumb_mini.webp"
		upload(localPath, cloudPath, "image/webp")
	}()

	_, err = firestoreClient.Collection(datafilesCollection).Doc(regionID).Set(context.Background(), datafileData)
	if err != nil {
		log.Fatalf("error updating document %#v in /datafiles: %v\n", regionID, err)
	}
}

func askForConfirmation() (bool, error) {
	var response string
	_, err := fmt.Scan(&response)
	if err != nil {
		return false, errors.WithStack(err)
	}
	if response == "y" || response == "Y" || response == "\n" {
		return true, nil
	} else if response == "N" || response == "n" {
		return false, nil
	}

	return false, nil
}

// Upload uploads file at localPath to Cloud Storage at cloudPath.
func upload(localPath string, cloudPath string, contentType string) {
	ctx := context.TODO()

	compressedDatafile, err := os.Open(localPath)
	if err != nil {
		log.Fatalln("upload: error opening compressed datafile:", err)
	}
	defer compressedDatafile.Close()

	bucket := storageClient.Bucket(bucketName)
	w := bucket.Object(cloudPath).NewWriter(ctx)
	w.ContentType = contentType
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	fmt.Printf("upload: begin uploading to %s...", cloudPath)
	_, err = io.Copy(w, compressedDatafile)
	if err != nil {
		log.Fatalln("upload: error copying compressedDatafile to writer:", err)
	}

	err = w.Close()
	if err != nil {
		log.Fatalln("upload: error closing storage writer:", err)
	}

	fmt.Println("ok")
}
