package upload

import "time"

// FirestoreDatafile represents a document in /datafiles collection in Firestore.
type FirestoreDatafile struct {
	Available        bool      `json:"available" firestore:"available"`
	Featured         []string  `json:"featured" firestore:"featured"`
	FileSize         int64     `json:"fileSize" firestore:"fileSize"`
	FileURL          string    `json:"fileURL" firestore:"fileURL"`
	LastUploadedTime time.Time `json:"lastUploadedTime" firestore:"lastUploadedTime"`
	GeneratedAt      time.Time `json:"generatedAt" firestore:"generatedAt"`
	UploadedAt       time.Time `json:"uploadedAt" firestore:"uploadedAt"`
	Position         int       `json:"position" firestore:"position"`
	RegionID         string    `json:"regionID" firestore:"regionID"`
	RegionName       string    `json:"regionName" firestore:"regionName"`
	IsTestVersion    bool      `json:"isTestVersion" firestore:"isTestVersion"`
	ThumbBlurhash    string    `json:"thumbBlurhash" firestore:"thumbBlurhash"`
	ThumbMiniURL     string    `json:"thumbMiniURL" firestore:"thumbMiniURL"`
	ThumbURL         string    `json:"thumbURL" firestore:"thumbURL"`
}
