// Package upload implements functionality related to uploading
// region's zip archive to the cloud.
package upload

import (
	"context"
	"encoding/json"
	"errors"
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
	firestoreClient *firestore.Client
	storageClient   *storage.Client
)

func init() {
	log.SetFlags(0)
}

func InitFirebase() error {
	opt := option.WithCredentialsFile("./key.json")

	var err error
	firestoreClient, err = firestore.NewClient(context.Background(), "discoverrudy", opt)
	if err != nil {
		return fmt.Errorf("initialize firestore: %v", err)
	}

	storageClient, err = storage.NewClient(context.Background(), opt)
	if err != nil {
		return fmt.Errorf("initialize storage: %v", err)
	}

	return nil
}

func Upload(regionID string, lang string, position int, onlyMeta bool, prod bool) error {
	zipFilePath := "compressed/" + regionID + ".zip"
	zipFileInfo, err := os.Stat(zipFilePath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("datafile archive %s doesn't exist", zipFilePath)
	}

	prefixedRegionID := regionID
	datafilesCollection := "datafiles"
	if !prod {
		prefixedRegionID += "Test"
		datafilesCollection += "Test"
	}

	// https://firebasestorage.googleapis.com/v0/b/discoverrudy.appspot.com/o/static %2Frudy%2Frudy.zip?alt=media
	fileLocation := appspotURL + url.QueryEscape("/"+prefixedRegionID+"/"+zipFileInfo.Name()) + "?alt=media"
	thumbLocation := appspotURL + url.QueryEscape("/"+prefixedRegionID+"/thumb.webp") + "?alt=media"
	thumbMiniLocation := appspotURL + url.QueryEscape("/"+prefixedRegionID+"/thumb_mini.webp") + "?alt=media"

	log.Println("fileLocation:", fileLocation)
	log.Println("thumbLocation:", thumbLocation)
	log.Println("thumbMiniLocation:", thumbMiniLocation)

	log.Println("making thumb blurhash...")
	thumbBlurhash, err := makeThumbBlurhash(regionID)
	if err != nil {
		return fmt.Errorf("make blurhash: %v", err)
	}

	fmt.Println("you are going to upload a data pack with the following metadata")
	meta, err := parseMeta(regionID, lang)
	if err != nil {
		return fmt.Errorf("parse meta: %v", err)
	}

	datafileData := FirestoreDatafile{
		Available:        true,
		Featured:         meta.Featured,
		FileSize:         zipFileInfo.Size(),
		FileURL:          fileLocation,
		LastUploadedTime: readers.CurrentTime(),
		GeneratedAt:      meta.GeneratedAt,
		UploadedAt:       readers.CurrentTime(),
		IsTestVersion:    !prod,
		Position:         position,
		RegionID:         regionID,
		RegionName:       meta.RegionName,
		ThumbBlurhash:    thumbBlurhash,
		ThumbMiniURL:     thumbMiniLocation,
		ThumbURL:         thumbLocation,
	}

	datafileDataJSON, err := json.MarshalIndent(datafileData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal datafileData to JSON: %v", err)
	}
	fmt.Println(string(datafileDataJSON))

	accepted, err := readers.AskForConfirmation(os.Stdin, os.Stdout, "upload: continue?", false)
	if err != nil {
		return fmt.Errorf("ask for confirmation: %v", err)
	}

	if !accepted {
		log.Println("operation canceled by the user")
		return nil
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
	log.Printf("updating document at %s...\n", docRef.Path)
	_, err = docRef.Set(context.Background(), datafileData)
	if err != nil {
		return fmt.Errorf("error updating document %#v in /datafiles: %v", regionID, err)
	}

	return nil
}

// Upload uploads file at localPath (relative) to Cloud Storage at cloudPath (absolute).
func upload(localPath string, cloudPath string, contentType string) error {
	ctx := context.TODO()

	compressedDatafile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open compressed datafile: %v", err)
	}
	defer compressedDatafile.Close()

	bucket := storageClient.Bucket(bucketName)
	w := bucket.Object(cloudPath).NewWriter(ctx)
	w.ContentType = contentType
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	fmt.Printf("uploading to %s...\n", cloudPath)

	_, err = io.Copy(w, compressedDatafile)
	if err != nil {
		return fmt.Errorf("copy compressedDatafile to writer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("close storage writer: %v", err)
	}

	return nil
}
