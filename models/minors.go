package models

// Action usually represents a URL with a name.
type Action struct {
	Name  Text   `json:"name"`
	Value string `json:"value"`
}

// Location represents single a point in the real world.
type Location struct {
	Lat float32 `json:"lat" firestore:"lat"`
	Lng float32 `json:"lng" firestore:"lng"`
}

type Link struct {
	Name       string `json:"name"`
	WebsiteURL string `json:"website_url"`
}

// Quality represents quality of an image.
type Quality int

const (
	// Compressed quality is most often used.
	Compressed = iota + 1
	// Original quality represents full, uncompressed image.
	Original
)
