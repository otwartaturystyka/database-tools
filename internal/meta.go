package internal

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/bartekpacia/database-tools/readers"
)

// Meta represents the JSON object in the beginning of data.json file.
type Meta struct {
	// Short, lowercase ID of the datafile's region.
	RegionID string `json:"region_id"`

	// Full localized name of the datafile's region.
	RegionName string `json:"region_name"`

	// Center of the Region
	Center Location `json:"center"`

	// Time of datafile generation. It is present only in generated datafile
	// i.e after the "generate" program has been run.
	GeneratedAt time.Time `json:"generated_at"`

	// People who somehow helped with creating the datafile.
	Contributors []string `json:"contributors"`

	// Some featured places present in the datafile.
	Featured []string `json:"featured"`

	// Resources (websites, books) which provided data in the datafile.
	Sources []Link `json:"sources"`

	// Related resources which might interest people using this datafile.
	Links []Link `json:"links"`
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
