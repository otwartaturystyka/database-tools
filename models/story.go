package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bartekpacia/database-tools/readers"
)

// Story represents a longer piece of text about a particular topic.
type Story struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MarkdownFile string `json:"markdown_filename"`
	markdownPath string
	Images       []string `json:"images"`
	imagePaths   []string
}

// Parse parses story data from its directory and assigns
// it to story pointed to by s. It must be used directly
// in the tracks's directory.
func (s *Story) Parse(lang string) error {
	name, err := readers.ReadFromFile(lang + "/name.txt")
	if err != nil {
		return err
	}
	s.Name = string(name)

	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		return err
	}

	err = s.makeMarkdownPath(lang)
	if err != nil {
		return fmt.Errorf("make markdown path: %w", err)
	}

	s.imagePaths = make([]string, 0)
	err = s.makeImagePaths(Compressed)
	if err != nil {
		return fmt.Errorf("make images paths: %w", err)
	}

	return nil
}

func (s *Story) makeMarkdownPath(lang string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	s.markdownPath = wd + "/" + lang + "/" + s.MarkdownFile + ".md"
	return nil
}

func (s *Story) makeImagePaths(quality Quality) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	var qualityDir string
	if quality == Compressed {
		qualityDir = "compressed"
	} else if quality == Original {
		qualityDir = "original"
	}

	// s.Images were set when the story was parsed from its JSON.
	if s.Images == nil {
		s.Images = make([]string, 0)
	}

	for _, image := range s.Images {
		absPath := filepath.Join(cwd, "images/", qualityDir, "/"+image+".webp")
		s.imagePaths = append(s.imagePaths, absPath)
	}

	return nil
}

// ImagePaths returns paths of all images of story s. They are
// specific to your machine!
func (s *Story) ImagePaths() []string {
	return s.imagePaths
}

// MarkdownPath returns path to the story's markdown file.
// It is specific to your machine.
func (s *Story) MarkdownPath() string {
	return s.markdownPath
}
