package internal

import (
	"encoding/json"
	"github.com/bartekpacia/database-tools/readers"
)

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
	name, err := readers.ReadFromFile("content/" + lang + "/name.txt")
	if err != nil {
		return err
	}
	dayroom.Name = string(name)

	overview, err := readers.ReadFromFile("content/" + lang + "/overview.txt")
	if err != nil {
		return err
	}
	dayroom.Overview = string(overview)

	quickInfo, err := readers.ReadFromFile("content/" + lang + "/quick_info.txt")
	if err != nil {
		return err
	}
	dayroom.QuickInfo = string(quickInfo)

	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &dayroom)
	if err != nil {
		return err
	}

	return nil
}

