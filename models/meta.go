package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/opentouristics/database-tools/readers"
)

// Meta represents the JSON object in the beginning of data.json file.
type Meta struct {
	// Short, lowercase ID of the datafile's region.
	RegionID string `json:"region_id"`

	// Full localized name of the datafile's region.
	RegionName Text `json:"region_name"`

	// Center of the Region
	Center Location `json:"center"`

	// Time of datafile generation. It is present only in generated datafile i.e
	// after the "generate" program has been run.
	GeneratedAt time.Time `json:"generated_at"`

	// People who somehow helped with creating the datafile.
	Contributors []string `json:"contributors"`

	// Some featured places present in the datafile.
	Featured []string `json:"featured"`

	// Resources (websites, books) which provided data in the datafile.
	Sources []Link `json:"sources"`

	// Related resources which might interest people using this datafile.
	Links []Link `json:"links"`

	// Hash that identifies the commit from which this datafile was generated.
	CommitHash string `json:"commit_hash"`

	// Non-nil only for production datafiles.
	CommitTag *string `json:"commit_tag"`

	// Count of places in all sections.
	PlaceCount int `json:"place_count"`
}

// Parse parses datafile's metadata and assigns it to meta struct pointed to by
// m.
func (m *Meta) Parse() error {
	name, err := readers.ReadLocalizedFiles("name.txt")
	if err != nil {
		return err
	}

	for k, v := range name {
		name[k] = strings.TrimSuffix(v, "\n")
	}
	m.RegionName = name

	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("unmarshal JSON: %v", err)
	}

	return nil
}

// ParseFromGenerated parses metadata of the datafile in the generated
// directory. It looks for data.json in the current dir, parses it and and
// assigns it to meta struct pointed to by m.
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
