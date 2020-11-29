package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Datafile represents structure of data.json file.
type Datafile struct {
	Meta     Meta      `json:"meta"`
	Sections []Section `json:"sections"`
	Tracks   []Track   `json:"tracks"`
	Stories  []Story   `json:"stories"`
	Dayrooms []Dayroom `json:"dayrooms"`
}

// Meta represents JSON object in the beginning of data.json file.
type Meta struct {
	RegionID     string              `json:"region_id"`
	RegionName   string              `json:"region_name"`
	GeneratedAt  string              `json:"generated_at"`
	Contributors []string            `json:"contributors"`
	Sources      []map[string]string `json:"sources"`
}

type Parseable interface {
	Parse(lang string) error
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
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Section   string   `json:"section"`
	Name      string   `json:"name"`
	QuickInfo string   `json:"quick_info"`
	Overview  string   `json:"overview"`
	Images    []string `json:"images"`
	Lat       float32  `json:"lat"`
	Lng       float32  `json:"lng"`
	Leader    string   `json:"leader"`
}

// Parse parses dayroom data from its directory and assigns
// it to dayroom pointer to by d. It must be used directly
// in the dayroom's directory, usually by using os.Chdir().
func (dayroom *Dayroom) Parse(lang string) error {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		return err
	}

	nameFile, err := os.Open("content/" + lang + "/name.txt")
	if err != nil {
		return err
	}

	overviewFile, err := os.Open("content/" + lang + "/overview.txt")
	if err != nil {
		return err
	}

	overview, err := ioutil.ReadAll(overviewFile)
	if err != nil {
		return err
	}
	dayroom.Overview = string(overview)

	quickInfoFile, err := os.Open("content/" + lang + "/quick_info.txt")
	if err != nil {
		return err
	}

	quickInfo, err := ioutil.ReadAll(quickInfoFile)
	if err != nil {
		return err
	}
	dayroom.QuickInfo = string(quickInfo)

	name, err := ioutil.ReadAll(nameFile)
	if err != nil {
		return err
	}
	dayroom.Name = string(name)

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &dayroom)
	if err != nil {
		return err
	}

	return nil
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
