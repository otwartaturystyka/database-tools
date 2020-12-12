package internal

import "time"

// FirestoreDatafile represents a document in /datafiles collection in Firestore.
type FirestoreDatafile struct {
	Available        bool      `firestore:"available"`
	Featured         []string  `firestore:"featured"`
	FileSize         int64     `firestore:"fileSize"`
	FileURL          string    `firestore:"fileURL"`
	LastUploadedTime time.Time `firestore:"lastUploadedTime,serverTimestamp"`
	Position         int       `firestore:"position"`
	RegionID         string    `firestore:"regionID"`
	RegionName       string    `firestore:"regionName"`
	IsTestVersion    bool      `firestore:"isTestVersion"`
	ThumbBlurhash    string    `firestore:"thumbBlurhash"`
	ThumbMiniURL     string    `firestore:"thumbMiniURL"`
	ThumbURL         string    `firestore:"thumbURL"`
}
