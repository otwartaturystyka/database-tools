package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bartekpacia/database-tools/readers"

	"github.com/pkg/errors"
)

// Datafile represents structure of data.json file.
type Datafile struct {
	Meta     Meta      `json:"meta"`
	Sections []Section `json:"sections"`
	Tracks   []Track   `json:"tracks"`
	Stories  []Story   `json:"stories"`
	Dayrooms []Dayroom `json:"dayrooms"`
}

// Meta represents the JSON object in the beginning of data.json file.
type Meta struct {
	// Short, lowercase ID of the datafile's region.
	RegionID string `json:"region_id"`

	// Full localized name of the datafile's region.
	RegionName string `json:"region_name"`

	// Time of datafile generation. It is present only in generated datafile
	// i.e after the "generate" program has been run.
	GeneratedAt time.Time `json:"generated_at"`

	// People who somehow helped with creating the datafile.
	Contributors []string `json:"contributors"`

	// Some featured places present in the datafile.
	Featured []string `json:"featured"`

	// Resources (websites, books) which provided data in the datafile.
	Sources []struct {
		Name       string `json:"name"`
		WebsiteURL string `json:"website_url"`
	} `json:"sources"`
}

// Parse parses datafile's metadata and assigns it to meta
// struct pointed to by m.
func (m *Meta) Parse(lang string) error {
	name, err := readers.ReadFromFile(filepath.Join(lang, "name.txt"))
	if err != nil {
		return err
	}
	m.RegionName = strings.TrimSuffix(string(name), "\n")

	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	return nil
}

// ParseFromGenerated parses metadata of the datafile in the
// generated directory. It looks for data.json in the current
// dir, parses it and and assigns it to meta struct pointed to by m.
func (m *Meta) ParseFromGenerated() error {
	datafileData, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	var datafile Datafile
	err = json.Unmarshal(datafileData, &datafile)
	if err != nil {
		return err
	}
	*m = datafile.Meta

	return nil
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

// Parse parses section data from its directory and assigns
// it to section pointed to by s. It must be used directly
// in the scetions's directory. It recursively parses places.
func (section *Section) Parse(lang string) error {
	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, section)
	if err != nil {
		return err
	}

	name, err := readers.ReadFromFile("content/" + lang + "/name.txt")
	if err != nil {
		return err
	}
	section.Name = string(name)

	quickInfo, err := readers.ReadFromFile("content/" + lang + "/quick_info.txt")
	if err != nil {
		return err
	}
	section.QuickInfo = string(quickInfo)

	// Parse places.
	places := make([]Place, 0, 50)
	placesWalker := func(path string, info os.FileInfo, err error) error {

		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}

		os.Chdir(path)

		var place Place
		err = place.Parse(lang)
		if err != nil {
			return errors.Wrapf(err, "parse place \"%s\"", path)
		}
		os.Chdir("../..")

		places = append(places, place)
		return nil
	}

	err = filepath.Walk("places", placesWalker)
	if err != nil {
		return errors.Wrap(err, "walk places")
	}

	section.Places = places

	return nil
}

// Quality represents the quality of the image.
type Quality int

const (
	// Compressed quality is most often used.
	Compressed = iota + 1
	// Original quality represents full, uncompressed image.
	Original
)
