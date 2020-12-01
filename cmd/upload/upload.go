package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/bartekpacia/database-tools/internal"

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
	position int
	verbose  bool
	test     bool
)

var (
	firestoreClient *firestore.Client
	storageClient   *storage.Client
)

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.IntVar(&position, "position", 1, "position at which the datafile should show in the app")
	flag.BoolVar(&test, "test", true, "whether to upload to test collections in Firestore")
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

	fmt.Printf("upload: begin making thumb blurhash...")
	thumbBlurhash, err := makeThumbBlurhash(regionID)
	if err != nil {
		log.Fatalln("upload: error making a blurhash:", err)
	}
	fmt.Println("ok")

	// Upload compressed datafile
	func() {
		localPath := "compressed/" + regionID + ".zip"
		cloudPath := "static/" + storagePrefix + "/rudy.zip"
		upload(localPath, cloudPath, "application/zip")
	}()

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

	// Upload minifed thumb

	meta := internal.FirestoreDatafile{
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

	_, err = firestoreClient.Collection(datafilesCollection).Doc(regionID).Set(context.Background(), meta)
	if err != nil {
		log.Fatalf("error updating document %#v in /datafiles: %v\n", regionID, err)
	}
}

// upload uploads a file under path to Cloud Storage path
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
