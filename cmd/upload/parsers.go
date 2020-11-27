package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// parseFeatured parses an array of featured for the datafile regionID from the database.
func parseFeatured(regionID string) (contributors []string, err error) {
	contribsFile, err := os.Open("database/" + regionID + "/meta/featured.json")
	if err != nil {
		// log.Fatalln("upload: error opening contributors file:", err)
		return nil, err
	}
	defer contribsFile.Close()

	b, err := ioutil.ReadAll(contribsFile)
	if err != nil {
		// log.Fatalln("upload: error reading from contributors file:", err)
		return nil, err
	}

	err = json.Unmarshal(b, &contributors)
	if err != nil {
		// log.Fatalln("upload: error unmarshalling contributors file:", err)
		return nil, err
	}

	return
}
