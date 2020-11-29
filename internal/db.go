package internal

// Meta represents JSON object in the beginning of data.json file.
type Meta struct {
	RegionID     string              `json:"region_id"`
	RegionName   string              `json:"region_name"`
	GeneratedAt  string              `json:"generated_at"`
	Contributors []string            `json:"contributors"`
	Sources      []map[string]string `json:"sources"`
}

// Section represents places of similiar type and associated metadata.
type Section struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	BgImage   string  `json:"background_image"`
	QuickInfo string  `json:"quick_info"`
	Places    []Place `json:"places"`
}

// Place represents single place in real world.
type Place struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Section     string   `json:"section"`
	Icon        string   `json:"icon"`
	QuickInfo   string   `json:"quick_info"`
	Overview    string   `json:"overview"`
	Lat         float32  `json:"lat"`
	Lng         float32  `json:"lng"`
	WebsiteURL  string   `json:"website_url"`
	FacebookURL string   `json:"facebook_url"`
	Headers     []string `json:"headers"`
	Content     []string `json:"content"`
	Images      []string `json:"images"`
}

// Track represents a bike trail.
type Track struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	QuickInfo string   `json:"quick_info"`
	Overview  string   `json:"overview"`
	Images    []string `json:"images"`
	Coords    []struct {
		Lat float32 `json:"lat"`
		Lng float32 `json:"lng"`
	} `json:"coords"`
}

// Story represents a longer piece of text about a particular topic.
type Story struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	MarkdownFile string   `json:"markdown_filename"`
	Images       []string `json:"images"`
}

// Dayroom represents a place run by local community.
type Dayroom struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	QuickInfo string `json:"quick_info"`
	Overview  string `json:"overview"`
	// Images    []string `json:"images"`
	Lat    float32 `json:"lat"`
	Lng    float32 `json:"lng"`
	Leader string  `json:"leader"`
}

// Quality represents the quality of the image.
type Quality int

const (
	// Compressed quality is most often used.
	Compressed = 1
	Original   = 2
)

// // Image represents a single image with some metadata *in the datafile*.
// type Image struct {
// 	Name    string
// 	Ext     string
// 	quality ImageQuality
// }
