package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"golang.org/x/image/webp"

	"github.com/bbrks/go-blurhash"
)

func makeThumbBlurhash(regionID string) (blur string, err error) {
	file, err := os.Open("database/" + regionID + "/meta/thumb.webp")
	if err != nil {
		return "", err
	}

	thumbImage, err := webp.Decode(file)
	if err != nil {
		return "", err
	}

	blur, err = blurhash.Encode(4, 3, thumbImage)
	if err != nil {
		return "", nil
	}

	return
}

// parseFeatured parses an array of featured for the datafile regionID from the database.
func parseFeatured(regionID string) (featured []string, err error) {
	contribsFile, err := os.Open("database/" + regionID + "/meta/featured.json")
	if err != nil {
		// log.Fatalln("upload: error opening featured file:", err)
		return nil, err
	}
	defer contribsFile.Close()

	b, err := ioutil.ReadAll(contribsFile)
	if err != nil {
		// log.Fatalln("upload: error reading from featured file:", err)
		return nil, err
	}

	err = json.Unmarshal(b, &featured)
	if err != nil {
		// log.Fatalln("upload: error unmarshalling featured file:", err)
		return nil, err
	}

	return
}
