package upload

import (
	"time"

	"github.com/opentouristics/database-tools/models"
)

// FirestoreDatafile represents a document in datafiles collection in Firestore.
type FirestoreDatafile struct {
	Available     bool              `json:"available" firestore:"available"`
	Featured      []string          `json:"featured" firestore:"featured"`
	FileSize      int64             `json:"fileSize" firestore:"fileSize"`
	FileURL       string            `json:"fileURL" firestore:"fileURL"`
	PlaceCount    int               `json:"placeCount" firestore:"placeCount"`
	GeneratedAt   time.Time         `json:"generatedAt" firestore:"generatedAt"`
	UploadedAt    time.Time         `json:"uploadedAt" firestore:"uploadedAt"`
	Position      int               `json:"position" firestore:"position"`
	RegionID      string            `json:"regionID" firestore:"regionID"`
	RegionName    models.Text       `json:"regionName" firestore:"regionName"`
	CommitHash    string            `json:"commitHash" firestore:"commitHash"`
	CommitTag     *string           `json:"commitTag" firestore:"commitTag"`
	IsTestVersion bool              `json:"isTestVersion" firestore:"isTestVersion"`
	ThumbBlurhash string            `json:"thumbBlurhash" firestore:"thumbBlurhash"`
	ThumbMiniURL  string            `json:"thumbMiniURL" firestore:"thumbMiniURL"`
	ThumbURL      string            `json:"thumbURL" firestore:"thumbURL"`
	Center        models.Location   `json:"center" firestore:"center"`
	Bounds        []models.Location `json:"bounds" firestore:"bounds"`
}
