package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/bartekpacia/database-tools/readers"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	bucketName = "discoverrudy.appspot.com"
	appspotURL = "https://firebasestorage.googleapis.com/v0/b/" + bucketName + "/o/static"
)

var (
	regionID        string
	lang            string
	displayPosition int
	onlyMeta        bool
	verbose         bool
)

var (
	firestoreClient *firestore.Client
	storageClient   *storage.Client
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.StringVar(&lang, "lang", "pl", "language of the datafile to upload")
	flag.IntVar(&displayPosition, "position", 1, "position at which the datafile will show in the app")
	flag.BoolVar(&onlyMeta, "only-meta", false, "upload only metadata (not the .zip file)")
	flag.BoolVar(&verbose, "verbose", false, "print extensive logs")

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
	zipFileInfo, err := os.Stat(zipFilePath)
	if os.IsNotExist(err) {
		log.Fatalf("upload: datafile archive %s doesn't exist\n", zipFilePath)
	}

	prefixedRegionID := regionID
	datafilesCollection := "datafiles"

	// https://firebasestorage.googleapis.com/v0/b/discoverrudy.appspot.com/o/static %2Frudy%2Frudy.zip?alt=media
	fileLocation := appspotURL + url.QueryEscape("/"+prefixedRegionID+"/"+zipFileInfo.Name()) + "?alt=media"
	thumbLocation := appspotURL + url.QueryEscape("/"+prefixedRegionID+"/thumb.webp") + "?alt=media"
	thumbMiniLocation := appspotURL + url.QueryEscape("/"+prefixedRegionID+"/thumb_mini.webp") + "?alt=media"

	fmt.Println("upload: fileLocation:", fileLocation)
	fmt.Println("upload: thumbLocation:", thumbLocation)
	fmt.Println("upload: thumbMiniLocation:", thumbMiniLocation)

	fmt.Printf("upload: making thumb blurhash...")
	thumbBlurhash, err := makeThumbBlurhash(regionID)
	if err != nil {
		log.Fatalf("\nupload: error making a blurhash: %v\n", err)
	}
	fmt.Println("ok")

	fmt.Println("upload: you are going to upload a data pack with the following metadata")

	meta, err := parseMeta(regionID, lang)
	if err != nil {
		log.Fatalln("upload: error parsing meta:", err)
	}

	datafileData := FirestoreDatafile{
		Available:        true,
		Featured:         meta.Featured,
		FileSize:         zipFileInfo.Size(),
		FileURL:          fileLocation,
		LastUploadedTime: readers.CurrentTime(),
		GeneratedAt:      meta.GeneratedAt,
		UploadedAt:       readers.CurrentTime(),
		IsTestVersion:    false,
		Position:         displayPosition,
		RegionID:         regionID,
		RegionName:       meta.RegionName,
		ThumbBlurhash:    thumbBlurhash,
		ThumbMiniURL:     thumbMiniLocation,
		ThumbURL:         thumbLocation,
	}

	datafileDataJSON, err := json.MarshalIndent(datafileData, "", "  ")
	if err != nil {
		log.Fatalln("upload: failed to marshal datafileData to JSON:", err)
	}
	fmt.Println(string(datafileDataJSON))

	accepted, err := readers.AskForConfirmation(os.Stdin, os.Stdout, "upload: continue?", false)
	if err != nil {
		log.Fatalf("\nupload: failed to get response: %v\n", err)
	}

	if !accepted {
		fmt.Println("upload: operation canceled by the user")
		os.Exit(0)
	}

	// Upload compressed datafile
	if !onlyMeta {
		func() {
			localPath := filepath.Join("compressed", regionID+".zip")
			cloudPath := path.Join("static", prefixedRegionID, regionID+".zip")
			upload(localPath, cloudPath, "application/zip")
		}()
	}

	// Upload thumb
	func() {
		localPath := filepath.Join("database", regionID+"/meta/thumb.webp")
		cloudPath := path.Join("static", prefixedRegionID, "thumb.webp")
		upload(localPath, cloudPath, "image/webp")
	}()

	// Upload minified thumb
	func() {
		localPath := filepath.Join("database", regionID+"/meta/thumb_mini.webp")
		cloudPath := path.Join("static", prefixedRegionID, "thumb_mini.webp")
		upload(localPath, cloudPath, "image/webp")
	}()

	docRef := firestoreClient.Collection(datafilesCollection).Doc(regionID)
	fmt.Printf("upload: updating document at %s...", docRef.Path)
	_, err = docRef.Set(context.Background(), datafileData)
	if err != nil {
		log.Fatalf("\nerror updating document %#v in /datafiles: %v\n", regionID, err)
	}
	fmt.Println("ok")
}

// Upload uploads file at localPath (relative) to Cloud Storage at cloudPath (absolute).
func upload(localPath string, cloudPath string, contentType string) {
	ctx := context.TODO()

	compressedDatafile, err := os.Open(localPath)
	if err != nil {
		log.Fatalf("\nupload: error opening compressed datafile: %v\n", err)
	}
	defer compressedDatafile.Close()

	bucket := storageClient.Bucket(bucketName)
	w := bucket.Object(cloudPath).NewWriter(ctx)
	w.ContentType = contentType
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	fmt.Printf("upload: uploading to %s...", cloudPath)
	_, err = io.Copy(w, compressedDatafile)
	if err != nil {
		log.Fatalf("\nupload: error copying compressedDatafile to writer: %v\n", err)
	}

	err = w.Close()
	if err != nil {
		log.Fatalf("\nupload: error closing storage writer: %v\n", err)
	}

	fmt.Println("ok")
}
