package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bartekpacia/database-tools/internal"
	"github.com/pkg/errors"
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
	if err != nil {
		log.Fatalln(errors.Unwrap(err))
	}
	datafile.Meta = meta

	sections, err := parseSections(lang)
	check(err)
	if err != nil {
		log.Fatalf("createOutputDir: %+v\n", err)
	}
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
	if err != nil {
		log.Fatalf("createOutputDir(): %v\n", err)
	}

	b, err := json.MarshalIndent(datafile, "", "	")
	check(err)

	n, err := dataJSONFile.Write(b)
	check(err)

	err = copyMarkdownFiles(&stories)
	check(err)

	// TODO err = copyImages(&images)

	fmt.Printf("generate: wrote %d KB to data.json file\n", n/1024)
}

func copyImages(places *[]internal.Place) error {
	return nil
}

// Must be run from this project's root dir.
func copyMarkdownFiles(stories *[]internal.Story) error {
	for _, story := range *stories {
		wd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to get working directory")
		}

		srcPath := wd + "/database/" + regionID + "/stories/" + story.ID + "/" + lang + "/" + story.MarkdownFile
		dstPath := wd + "/generated/" + regionID + "/stories/" + story.MarkdownFile

		src, err := os.Open(srcPath)
		if err != nil {
			return errors.Errorf("failed to open markdown file at %s", srcPath)
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return errors.Errorf("failed to create dst file at %s", dstPath)
		}

		n, err := io.Copy(dst, src)
		check(err)

		fmt.Printf("generate: copied story %s, %.1f KB of text copied, ", story.ID, float32(n)/1000)
		if len(story.Images) > 0 {
			fmt.Printf("has %d images\n", len(story.Images))
		} else {
			fmt.Printf("has no images\n")
		}
	}

	return nil
}

// CreateOutputDir creates a datafile directory structure inside generated/ in project root.
func createOutputDir(regionID string) (*os.File, error) {
	generatedPath := "generated"
	outputDirPath := generatedPath + "/" + regionID

	if _, err := os.Stat(generatedPath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Errorf("dir %#v does not exist", generatedPath)
		}

		return nil, errors.Errorf("failed to stat %#v dir", generatedPath)
	}

	err := os.RemoveAll(outputDirPath)
	if err != nil {
		return nil, errors.Errorf("failed to remove output dir %#v", outputDirPath)
	}

	err = os.Mkdir(outputDirPath, 0755)
	if err != nil {
		return nil, errors.Errorf("failed to make dir %#v", outputDirPath)
	}

	imagesDirPath := outputDirPath + "/images"
	err = os.Mkdir(imagesDirPath, 0755)
	if err != nil {
		return nil, errors.Errorf("failed to make dir %#v (for images)", imagesDirPath)
	}

	storiesDirPath := outputDirPath + "/stories"
	err = os.Mkdir(storiesDirPath, 0755)
	if err != nil {
		return nil, errors.Errorf("failed to make dir %#v (for stories)", storiesDirPath)
	}

	dataJSONPath := outputDirPath + "/data.json"
	dataJSONFile, err := os.Create(dataJSONPath)
	if err != nil {
		return nil, errors.Errorf("failed to create file %#v (the main json file)", dataJSONPath)
	}

	return dataJSONFile, nil
}
