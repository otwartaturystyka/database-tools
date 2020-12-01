package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bartekpacia/database-tools/internal"
)

var (
	regionID string
	language string
	quality  int
	verbose  bool
)

func check(err error) {
	if err != nil {
		log.Fatalln("generate:", err)
	}
}

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.StringVar(&language, "lang", "pl", "language of text in the datafile")
	flag.IntVar(&quality, "quality", 1, "quality of photos in the datafile")
}

func main() {
	flag.Parse()

	var datafile internal.Datafile

	if regionID == "" {
		log.Fatalln("generate: regionID is empty")
	}

	if quality != 1 && quality != 2 {
		log.Fatalln("generate: quality is not 1 or 2")
	}

	err := os.Chdir("database/" + regionID)
	if err != nil {
		log.Fatalln("generate:", err)
	}

	meta, err := parseMeta(language)
	check(err)
	datafile.Meta = meta

	sections, err := parseSections(language)
	check(err)
	datafile.Sections = sections

	tracks, err := parseTracks(language)
	check(err)
	datafile.Tracks = tracks

	stories, err := parseStories(language)
	check(err)
	datafile.Stories = stories

	dayrooms, err := parseDayrooms(language)
	check(err)
	datafile.Dayrooms = dayrooms

	os.Chdir("../..")

	dataJSONFile, err := createOutputDir(regionID)
	check(err)

	b, err := json.MarshalIndent(datafile, "", "	")
	check(err)

	n, err := dataJSONFile.Write(b)
	check(err)

	fmt.Printf("generate: wrote %d KB to data.json file\n", n/1024)
}

// CreateOutputDir creates a datafile directory structure inside generated/ in project root.
func createOutputDir(regionID string) (*os.File, error) {
	outputDirPath := "generated/" + regionID

	err := os.RemoveAll(outputDirPath)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(outputDirPath, 0755)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(outputDirPath+"/images", 0755)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(outputDirPath+"/stories", 0755)
	if err != nil {
		return nil, err
	}

	dataJSONFile, err := os.Create("generated/" + regionID + "/data.json")
	if err != nil {
		return nil, err
	}

	return dataJSONFile, nil
}
