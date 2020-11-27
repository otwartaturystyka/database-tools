package main

import "time"

// Meta represents JSON object in the beginning of data.json file.
type Meta struct {
	RegionID     string              `json:"region_id"`
	RegionName   string              `json:"region_name"`
	GeneratedAt  string              `json:"generated_at"`
	Contributors []string            `json:"contributors"`
	Sources      []map[string]string `json:"sources"`
}

// Datafile represents a document in /datafiles collection in Firestore.
type Datafile struct {
	Available        bool      `firestore:"available,omitempty"`
	Featured         []string  `firestore:"featured,omitempty"`
	FileSize         int       `firestore:"fileSize,omitempty"`
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
