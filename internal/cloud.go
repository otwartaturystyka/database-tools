package internal

import "time"

// FirestoreDatafile represents a document in /datafiles collection in Firestore.
type FirestoreDatafile struct {
	Available     bool      `json:"available" firestore:"available"`
	Featured      []string  `json:"featured" firestore:"featured"`
	FileSize      int64     `json:"fileSize" firestore:"fileSize"`
	FileURL       string    `json:"fileURL" firestore:"fileURL"`
	GeneratedAt   time.Time `json:"generatedAt" firestore:"generatedAt,serverTimestamp"`
	UploadedAt    time.Time `json:"uploadedAt" firestore:"uploadedAt,serverTimestamp"`
	Position      int       `json:"position" firestore:"position"`
	RegionID      string    `json:"regionID" firestore:"regionID"`
	RegionName    string    `json:"regionName" firestore:"regionName"`
	IsTestVersion bool      `json:"isTestVersion" firestore:"isTestVersion"`
	ThumbBlurhash string    `json:"thumbBlurhash" firestore:"thumbBlurhash"`
	ThumbMiniURL  string    `json:"thumbMiniURL" firestore:"thumbMiniURL"`
	ThumbURL      string    `json:"thumbURL" firestore:"thumbURL"`
}

// CurrentTime returns current UTC time rounded to second.
// This method's purpose is to standarize time formats in this project.
// Use it.
func CurrentTime() MyTime {
	return MyTime{time.Now().Round(time.Second).UTC()}
}
