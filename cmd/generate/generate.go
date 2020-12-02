package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bartekpacia/database-tools/internal"
)

var (
	regionID string
	lang     string
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
	flag.StringVar(&lang, "lang", "pl", "language of text in the datafile")
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

	meta, err := parseMeta(lang)
	check(err)
	datafile.Meta = meta

	sections, err := parseSections(lang)
	check(err)
	datafile.Sections = sections

	tracks, err := parseTracks(lang)
	check(err)
	datafile.Tracks = tracks

	stories, err := parseStories(lang)
	check(err)
	datafile.Stories = stories

	dayrooms, err := parseDayrooms(lang)
	check(err)
	datafile.Dayrooms = dayrooms

	os.Chdir("../..")

	dataJSONFile, err := createOutputDir(regionID)
	check(err)

	b, err := json.MarshalIndent(datafile, "", "	")
	check(err)

	n, err := dataJSONFile.Write(b)
	check(err)

	err = copyMarkdownFiles(&stories)
	check(err)

	fmt.Printf("generate: wrote %d KB to data.json file\n", n/1024)
}

// Must be run from this project's root dir.
func copyMarkdownFiles(stories *[]internal.Story) error {
	for _, story := range *stories {
		wd, err := os.Getwd()
		check(err)

		srcPath := wd + "/database/" + regionID + "/stories/" + story.ID + "/" + lang + "/" + story.MarkdownFile
		dstPath := wd + "/generated/" + regionID + "/stories/" + story.MarkdownFile

		src, err := os.Open(srcPath)
		check(err)

		dst, err := os.Create(dstPath)
		check(err)

		n, err := io.Copy(dst, src)
		check(err)

		fmt.Printf("generate: copied story %s, %d bytes copied\n", story.ID, n)
	}

	return nil
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
