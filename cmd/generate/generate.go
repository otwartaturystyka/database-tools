package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bartekpacia/database-tools/readers"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/internal"
	"github.com/pkg/errors"
)

var (
	regionID string
	lang     string
	quality  int
	verbose  bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.StringVar(&lang, "lang", "pl", "language of text in the datafile")
	flag.IntVar(&quality, "quality", 1, "quality of photos in the datafile")
	flag.BoolVar(&verbose, "verbose", false, "print extensive logs")
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
	datafile.Meta.GeneratedAt = readers.CurrentTime() // Important!

	sections, err := parseSections(lang)
	if err != nil {
		log.Fatalf("generate: failed to parse sections: %v\n", err)
	}
	datafile.Sections = sections

	tracks, err := parseTracks(lang)
	if err != nil {
		log.Fatalf("generate: failed to parse tracks: %v\n", err)
	}
	datafile.Tracks = tracks

	stories, err := parseStories(lang)
	if err != nil {
		log.Fatalf("generate: parseStories(): %v\n", err)
	}
	datafile.Stories = stories

	dayrooms, err := parseDayrooms(lang)
	if err != nil {
		log.Fatalf("generate: failed to parse dayrooms: %v\n", err)
	}
	datafile.Dayrooms = dayrooms

	os.Chdir("../..")

	fmt.Printf("generate: creating output dir...")
	dataJSONFile, err := createOutputDir(regionID)
	if err != nil {
		log.Fatalf("\ngenerate: failed to create output dir: %v\n", err)
	}
	fmt.Println("ok")

	fmt.Printf("generate: marshalling datafile to JSON...")
	data, err := json.MarshalIndent(datafile, "", "	")
	if err != nil {
		log.Fatalf("\ngenerate: failed to marshal datafile to JSON: %v\n", err)
	}
	fmt.Println("ok")

	fmt.Printf("generate: writing datafile json to a file...")
	n, err := dataJSONFile.Write(data)
	if err != nil {
		log.Fatalf("\ngenerate: failed to write data to the JSON file: %v\n", err)
	}
	fmt.Println("ok")
	fmt.Printf("generate: wrote %d KB to data.json file\n", n/1024)

	for _, section := range sections {
		for _, place := range section.Places {
			for _, imagePath := range place.ImagePaths() {
				_, err = copyImage(regionID, imagePath)
				if err != nil {
					log.Fatalf("generate: %v\n", err)
				}

				if strings.HasPrefix(filepath.Base(imagePath), "ic_") {
					makeMiniIcon(imagePath)
				}
			}
		}
	}

	for _, story := range stories {
		_, err := copyMarkdown(regionID, story.MarkdownPath())
		if err != nil {
			log.Fatalf("generate: failed to copy markdown file for story %s: %v\n", story.ID, err)
		}

		for _, path := range story.ImagePaths() {
			_, err := copyImage(regionID, path)
			if err != nil {
				log.Fatalf("generate: failed to copy image for story %s: %v\n", story.ID, err)
			}
		}
	}
}

func makeMiniIcon(srcPath string) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get working dir")
	}

	miniIconFilename := "mini_" + filepath.Base(srcPath)
	dstPath := filepath.Join(wd, "generated", regionID, "images", miniIconFilename)
	// fmt.Printf("srcPath: %s, dstPath: %s\n", srcPath, dstPath)
	err = exec.Command("magick", srcPath, "-quality", "60%", "-resize", "128x128", dstPath).Run()
	if err != nil {
		return errors.Wrapf(err, "failed to make a mini icon of image at %s", srcPath)
	}

	return nil
}

func copyImage(regionID string, srcPath string) (int, error) {
	n, err := copyFile(regionID, srcPath, "images")
	return n, err
}

func copyMarkdown(regionID string, srcPath string) (int, error) {
	n, err := copyFile(regionID, srcPath, "stories")
	return n, err
}

func copyFile(regionID string, srcPath string, subdir string) (int, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get working dir")
	}

	dstPath := wd + "/generated/" + regionID + "/" + subdir + "/" + filepath.Base(srcPath)
	src, err := os.Open(srcPath)
	if err != nil {
		return 0, errors.Errorf("failed to open src file at %s", srcPath)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return 0, errors.Errorf("failed to create dst file at %s", dstPath)
	}

	n, err := io.Copy(dst, src)
	if err != nil {
		return 0, errors.Errorf("failed to copy file from %s to %s", srcPath, dstPath)
	}

	return int(n), nil
}

// CreateOutputDir creates a datafile directory structure inside generated/ in project root.
func createOutputDir(regionID string) (*os.File, error) {
	generatedPath := "generated"
	outputDirPath := generatedPath + "/" + regionID

	// Check if the generated dir exists...
	if _, err := os.Stat(generatedPath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(generatedPath, 0755)
			if err != nil {
				return nil, errors.Errorf("dir %#v does not exist and cannot be created", generatedPath)
			}
		} else {
			return nil, errors.Errorf("failed to stat %#v dir", generatedPath)
		}
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
