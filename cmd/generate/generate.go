// Package generate implements the process of generating a directory
// from data available in the database.
package generate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/models"
	"github.com/bartekpacia/database-tools/readers"
)

func init() {
	log.SetFlags(0)
}

// Generate walks the database and copies files from it to the generated directory.
func Generate(regionID string, lang string, quality models.Quality, verbose bool) error {
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

	meta, err := parseMeta(lang)
	if err != nil {
		return fmt.Errorf("parse meta: %v", err)
	}
	datafile.Meta = meta
	datafile.Meta.GeneratedAt = readers.CurrentTime() // Important!

	sections, err := parseSections(lang)
	if err != nil {
		return fmt.Errorf("parse sections: %v", err)
	}
	datafile.Sections = sections

	tracks, err := parseTracks(lang)
	if err != nil {
		return fmt.Errorf("parse tracks: %v", err)
	}
	datafile.Tracks = tracks

	stories, err := parseStories(lang)
	if err != nil {
		return fmt.Errorf("parse stories: %v", err)
	}
	datafile.Stories = stories

	os.Chdir("../..")

	log.Println("creating output dir...")
	dataJSONFile, err := createOutputDir(regionID)
	if err != nil {
		return fmt.Errorf("create output directory: %v", err)
	}

	log.Println("marshalling datafile to JSON...")
	data, err := json.MarshalIndent(datafile, "", "	")
	if err != nil {
		return fmt.Errorf("marshal datafile struct to JSON: %v", err)
	}

	log.Println("writing datafile json to a file...")
	n, err := dataJSONFile.Write(data)
	if err != nil {
		return fmt.Errorf("write data to JSON file: %v", err)
	}

	log.Printf("wrote %d KB to data.json file\n", n/1024)

	for _, section := range sections {
		for _, place := range section.Places {
			for _, imagePath := range place.ImagePaths() {
				_, err = copyImage(regionID, imagePath)
				if err != nil {
					return fmt.Errorf("copy image: %v", err)
				}

				if strings.HasPrefix(filepath.Base(imagePath), "ic_") {
					err = makeMiniIcon(regionID, imagePath)
					if err != nil {
						return fmt.Errorf("make mini icon at %s: %v", imagePath, err)
					}
				}
			}
		}
	}

	for _, story := range stories {
		_, err := copyMarkdown(regionID, story.MarkdownPath())
		if err != nil {
			return fmt.Errorf("copy markdown file for story %s: %v", story.ID, err)
		}

		for _, path := range story.ImagePaths() {
			_, err := copyImage(regionID, path)
			if err != nil {
				return fmt.Errorf("copy image for story %s: %v", story.ID, err)
			}
		}
	}

	return nil
}

func makeMiniIcon(regionID string, srcPath string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	miniIconFilename := "mini_" + filepath.Base(srcPath)
	dstPath := filepath.Join(wd, "generated", regionID, "images", miniIconFilename)
	err = exec.Command("magick", srcPath, "-quality", "60%", "-resize", "128x128", dstPath).Run()
	if err != nil {
		return fmt.Errorf("run ImageMagick on image at %s: %w", srcPath, err)
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
		return 0, fmt.Errorf("open src file at %s: %w", srcPath, err)
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

// CreateOutputDir creates a datafile directory structure inside generated/ in project root.
func createOutputDir(regionID string) (*os.File, error) {
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

	dataJSONPath := filepath.Join(outputDirPath, "data.json")
	dataJSONFile, err := os.Create(dataJSONPath)
	if err != nil {
		return nil, fmt.Errorf("create file %#v (the main json file): %w", dataJSONPath, err)
	}

	return dataJSONFile, nil
}
