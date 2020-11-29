package internal

import "time"

// Datafile represents a document in /datafiles collection in Firestore.
type Datafile struct {
	Available        bool      `firestore:"available,omitempty"`
	Featured         []string  `firestore:"featured,omitempty"`
	FileSize         int64     `firestore:"fileSize,omitempty"`
	FileURL          string    `firestore:"fileURL,omitempty"`
	LastUploadedTime time.Time `firestore:"time,serverTimestamp"`
	Position         int       `firestore:"position,omitempty"`
	RegionID         string    `firestore:"regionID,omitempty"`
	RegionName       string    `firestore:"regionName,omitempty"`
	IsTestVersion    bool      `firestore:"isTestVersion,omitempty"`
	ThumbBlurhash    string    `firestore:"thumbBlurhash,omitempty"`
	ThumbMiniURL     string    `firestore:"thumbMiniURL,omitempty"`
	ThumbURL         string    `firestore:"thumbURL,omitempty"`
}