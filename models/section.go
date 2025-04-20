package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opentouristics/database-tools/formatters"
	"github.com/opentouristics/database-tools/readers"
)

// Section represents places of similiar type and associated metadata.
type Section struct {
	ID        string  `json:"id"`
	Name      Text    `json:"name"`
	BgImage   string  `json:"background_image"`
	QuickInfo Text    `json:"quick_info"`
	Places    []Place `json:"places"`
}

// Parse parses section data from its directory and assigns
// it to section pointed to by s. It must be used directly
// in the scetions's directory. It recursively parses places.
func (section *Section) Parse(verbose bool) error {
	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, section)
	if err != nil {
		return err
	}

	name, err := readers.ReadLocalizedFiles("name.txt")
	if err != nil {
		return err
	}
	section.Name = formatters.ToContent(name)

	quickInfo, err := readers.ReadLocalizedFiles("quick_info.txt")
	if err != nil {
		return err
	}
	section.QuickInfo = formatters.ToContent(quickInfo)

	// Parse places.
	places := make([]Place, 0, 50)
	placesWalker := func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("start walk place %#v: %w", path, err)
		}

		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}

		os.Chdir(path)

		var place Place
		err = place.Parse(verbose)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
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
