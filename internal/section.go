package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/readers"
)

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

	name, err := readers.ReadFromFile(filepath.Join("content", lang, "name.txt"))
	if err != nil {
		return err
	}
	section.Name = string(name)

	quickInfo, err := readers.ReadFromFile(filepath.Join("content", lang, "quick_info.txt"))
	if err != nil {
		return err
	}
	section.QuickInfo = string(quickInfo)

	// Parse places.
	places := make([]Place, 0, 50)
	placesWalker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("parse place %#v: %w", path, err)
		}

		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}

		os.Chdir(path)

		var place Place
		err = place.Parse(lang)
		if err != nil {
			return fmt.Errorf("parse place %#v: %w", path, err)
		}
		os.Chdir("../..")

		places = append(places, place)
		return nil
	}

	err = filepath.Walk("places", placesWalker)
	if err != nil {
		return fmt.Errorf("walk places: %w", err)
	}

	section.Places = places

	return nil
}
