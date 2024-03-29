// Package generate implements the process of generating a directory from data
// available in the database.
package generate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/opentouristics/database-tools/models"
	"github.com/opentouristics/database-tools/readers"
)

func init() {
	log.SetFlags(0)
}

// Generate walks the database and copies files from it to the generated
// directory.
func Generate(regionID string, quality models.Quality, verbose bool) error {
	var datafile models.Datafile

	if regionID == "" {
		return fmt.Errorf("regionID is empty")
	}

	if quality != 1 && quality != 2 {
		return fmt.Errorf("quality is not 1 or 2")
	}

	err := os.Chdir(filepath.Join("datafiles", "datafile-"+regionID))
	if err != nil {
		return fmt.Errorf("chdir into datafile's directory: %v", err)
	}

	meta, err := parseMeta()
	if err != nil {
		return fmt.Errorf("failed to parse meta: %v", err)
	}
	datafile.Meta = meta
	datafile.Meta.GeneratedAt = readers.CurrentTime() // Important!

	sections, err := parseSections(verbose)
	if err != nil {
		return fmt.Errorf("failed to parse sections: %v", err)
	}
	datafile.Sections = sections
	datafile.Meta.PlaceCount = len(datafile.AllPlaces())

	tracks, err := parseTracks()
	if err != nil {
		return fmt.Errorf("failed to parse tracks: %v", err)
	}
	datafile.Tracks = tracks

	stories, err := parseStories()
	if err != nil {
		return fmt.Errorf("failed to parse stories: %v", err)
	}
	datafile.Stories = stories

	os.Chdir("../..")

	log.Println("creating output dir...")
	outputDirPath, err := createOutputDir(regionID)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	log.Println("marshalling datafile to JSON...")
	data, err := json.MarshalIndent(datafile, "", "	")
	if err != nil {
		return fmt.Errorf("failed to marshal datafile struct to JSON: %v", err)
	}

	log.Println("creating file for datafile JSON contents...")
	dataJSONFile, err := os.Create(filepath.Join(*outputDirPath, "data.json"))
	if err != nil {
		return fmt.Errorf("create dataJSONFile: %w", err)
	}

	log.Println("writing datafile JSON to a file...")
	n, err := dataJSONFile.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to JSON file: %v", err)
	}

	log.Printf("wrote %d KB to data.json file\n", n/1024)

	log.Println("marshalling meta to JSON...")
	data, err = json.MarshalIndent(meta, "", "	")
	if err != nil {
		return fmt.Errorf("failed to marshal datafile struct to JSON: %v", err)
	}

	log.Println("creating file for meta JSON contents...")
	metaJSONFile, err := os.Create(filepath.Join(*outputDirPath, "meta.json"))
	if err != nil {
		return fmt.Errorf("create metaJSONFile: %w", err)
	}

	log.Println("writing meta JSON to file")
	n, err = metaJSONFile.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write meta JSON file: %v", err)
	}

	log.Printf("wrote %d KB to meta.json file\n", n/1024)

	for _, section := range sections {
		for _, place := range section.Places {
			for _, imagePath := range place.ImagePaths() {
				_, err = copyImage(regionID, imagePath)
				if err != nil {
					return fmt.Errorf("failed to copy image: %v", err)
				}
			}
		}
	}

	for _, story := range stories {
		_, err := copyMarkdown(regionID, story.MarkdownPath())
		if err != nil {
			return fmt.Errorf("failed to copy markdown file for story %s: %v", story.ID, err)
		}

		for _, path := range story.ImagePaths() {
			_, err := copyImage(regionID, path)
			if err != nil {
				return fmt.Errorf("failed to copy image for story %s: %v", story.ID, err)
			}
		}
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
		return 0, fmt.Errorf("get working dir: %w", err)
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}

	dstPath := filepath.Join(wd, "generated", regionID, subdir, filepath.Base(srcPath))
	dst, err := os.Create(dstPath)
	if err != nil {
		return 0, fmt.Errorf("create dst file at %s: %w", dstPath, err)
	}

	n, err := io.Copy(dst, src)
	if err != nil {
		return 0, fmt.Errorf("copy file from %s to %s: %w", srcPath, dstPath, err)
	}

	return int(n), nil
}

// CreateOutputDir creates a datafile directory structure inside generated/ in
// project root.
func createOutputDir(regionID string) (*string, error) {
	generatedPath := "generated"
	outputDirPath := filepath.Join(generatedPath, regionID)

	// Check if the generated dir exists...
	_, err := os.Stat(generatedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(generatedPath, 0o755)
			if err != nil {
				return nil, fmt.Errorf("dir %#v does not exist and cannot be created: %w", generatedPath, err)
			}
		} else {
			return nil, fmt.Errorf("stat %#v dir: %w", generatedPath, err)
		}
	}

	err = os.RemoveAll(outputDirPath)
	if err != nil {
		return nil, fmt.Errorf("remove output dir %#v: %w", outputDirPath, err)
	}

	err = os.Mkdir(outputDirPath, 0o755)
	if err != nil {
		return nil, fmt.Errorf("make output dir %#v: %w", outputDirPath, err)
	}

	imagesDirPath := filepath.Join(outputDirPath, "images")
	err = os.Mkdir(imagesDirPath, 0o755)
	if err != nil {
		return nil, fmt.Errorf("make dir %#v (for images): %w", imagesDirPath, err)
	}

	storiesDirPath := filepath.Join(outputDirPath, "stories")
	err = os.Mkdir(storiesDirPath, 0o755)
	if err != nil {
		return nil, fmt.Errorf("make dir %#v (for stories): %w", storiesDirPath, err)
	}

	return &outputDirPath, nil
}
