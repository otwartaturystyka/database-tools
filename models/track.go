package models

import (
	"encoding/json"

	"github.com/opentouristics/database-tools/formatters"
	"github.com/opentouristics/database-tools/readers"
)

// Track represents a bike trail or some other "long" geographical object.
type Track struct {
	ID        string     `json:"id"`
	Name      Text       `json:"name"`
	QuickInfo Text       `json:"quick_info"`
	Overview  Text       `json:"overview"`
	Images    []string   `json:"images"`
	Coords    []Location `json:"coords"`
}

// Parse parses track data from its directory and assigns
// it to track pointed to by t. It must be used directly
// in the track's directory, usually by using os.Chdir().
func (t *Track) Parse() error {
	name, err := readers.ReadLocalizedFiles("name.txt")
	if err != nil {
		return err
	}
	t.Name = name

	overview, err := readers.ReadLocalizedFiles("overview.txt")
	if err != nil {
		return err
	}
	t.Overview = formatters.ToContent(overview)

	quickInfo, err := readers.ReadLocalizedFiles("quick_info.txt")
	if err != nil {
		return err
	}
	t.QuickInfo = formatters.ToContent(quickInfo)

	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, t)

	return err
}
