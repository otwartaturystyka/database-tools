package models

import (
	"encoding/json"

	"github.com/bartekpacia/database-tools/readers"
)

// Track represents a bike trail or some other "long" geographical object.
type Track struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	QuickInfo string     `json:"quick_info"`
	Overview  string     `json:"overview"`
	Images    []string   `json:"images"`
	Coords    []Location `json:"coords"`
}

// Parse parses track data from its directory and assigns
// it to track pointed to by t. It must be used directly
// in the track's directory, usually by using os.Chdir().
func (t *Track) Parse(lang string) error {
	name, err := readers.ReadFromFile(lang + "/name.txt")
	if err != nil {
		return err
	}
	t.Name = string(name)

	overview, err := readers.ReadFromFile(lang + "/overview.txt")
	if err != nil {
		return err
	}
	t.Overview = string(overview)

	quickInfo, err := readers.ReadFromFile(lang + "/quick_info.txt")
	if err != nil {
		return err
	}
	t.QuickInfo = string(quickInfo)

	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, t)

	return err
}
