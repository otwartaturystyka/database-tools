package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opentouristics/database-tools/formatters"
	"github.com/opentouristics/database-tools/readers"
)

// Place represents single place in real world.
type Place struct {
	ID          string   `json:"id"`
	Name        Text     `json:"name"`
	Section     string   `json:"section"`
	Icon        string   `json:"icon"`
	QuickInfo   Text     `json:"quick_info"`
	Overview    Text     `json:"overview"`
	Lat         float32  `json:"lat"`
	Lng         float32  `json:"lng"`
	WebsiteURL  *string  `json:"website_url"`
	FacebookURL *string  `json:"facebook_url"`
	Headers     []Text   `json:"headers"`
	Content     []Text   `json:"content"`
	Actions     []Action `json:"actions"`
	Images      []string `json:"images"`
	imagePaths  []string
}

// Parse parses place data from its directory and assigns
// it to track pointed to by p. It must be used directly
// in the place's directory.
func (p *Place) Parse(verbose bool) error {
	// Technical metadata
	data, err := readers.ReadFromFile("data.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	err = p.makeImagePaths(Compressed)
	if err != nil {
		return fmt.Errorf("make image paths for place %s: %w", p.ID, err)
	}

	// Content
	name, err := readers.ReadLocalizedFiles("name.txt")
	if err != nil {
		return fmt.Errorf("read localized name: %v", err)
	}
	p.Name = formatters.ToContent(name)

	quickInfo, err := readers.ReadLocalizedFiles("quick_info.txt")
	if err != nil {
		return fmt.Errorf("read localized quick info: %v", err)
	}
	p.QuickInfo = formatters.ToContent(quickInfo)

	overview, err := readers.ReadLocalizedFiles("overview.txt")
	if err != nil {
		return fmt.Errorf("read localized overview: %v", err)
	}
	p.Overview = formatters.ToContent(overview)

	// Headers and content
	p.Headers = make([]Text, 0)
	p.Content = make([]Text, 0)

	entries, err := os.ReadDir("content/pl")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	textFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), "text_") {
			textFiles = append(textFiles, entry.Name())
		}
	}

	for _, textFile := range textFiles {
		text, err := readers.ReadLocalizedFiles(textFile)
		if err != nil {
			return err
		}

		header, content := formatters.ToSection(text)

		p.Headers = append(p.Headers, header)
		p.Content = append(p.Content, content)
	}

	// Actions
	p.Actions = make([]Action, 0)
	err = p.makeActions(verbose)
	if err != nil {
		return fmt.Errorf("make actions for place %s: %w", p.ID, err)
	}

	return nil
}

func (p *Place) makeActions(verbose bool) error {
	actionValuesFile, err := readers.ReadFromFile("actions.json")
	if err != nil {
		if verbose {
			fmt.Printf(
				"failed to open file %s of place %s (most "+
					"probably, this place does not have any actions, so the file "+
					"does not not exist)\n",
				"actions.json",
				p.ID,
			)
		}

		return nil
	}

	// Read action values from JSON
	actionValues := make([]string, 0)
	err = json.Unmarshal(actionValuesFile, &actionValues)
	if err != nil {
		return err
	}

	// Read action name for every available language
	entries, err := os.ReadDir("content/pl")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	actionNameFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasPrefix(entry.Name(), "action_") {
			actionNameFiles = append(actionNameFiles, entry.Name())
		}
	}

	actionNames := make([]Text, 0)
	for _, actionNameFile := range actionNameFiles {
		text, err := readers.ReadLocalizedFiles(actionNameFile)
		if err != nil {
			return err
		}

		for key, value := range text {
			text[key] = strings.TrimSuffix(value, "\n")
		}

		actionNames = append(actionNames, text)
	}

	if len(actionValues) != len(actionNames) {
		return errors.New("actionValues and actionNames are not of the same length – this is probably an error in the database")
	}

	for i := 0; i < len(actionValues); i++ {
		action := Action{Name: actionNames[i], Value: actionValues[i]}
		p.Actions = append(p.Actions, action)
	}

	return nil
}

func (p *Place) makeImagePaths(quality Quality) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working dir: %w", err)
	}

	var qualityDir string
	if quality == Compressed {
		qualityDir = "compressed"
	} else if quality == Original {
		qualityDir = "original"
	}

	// s.Images were set when the story was parsed from its JSON
	for _, image := range p.Images {
		absImagePath := filepath.Join(workingDir, "images", qualityDir, image+".webp")

		if _, err := os.Stat(absImagePath); err != nil {
			return fmt.Errorf("image at %s does not exist: %w", absImagePath, err)
		}

		p.imagePaths = append(p.imagePaths, absImagePath)
	}

	// Add icon
	absIconPath := filepath.Join(workingDir, "images", qualityDir, p.Icon+".webp")
	p.imagePaths = append(p.imagePaths, absIconPath)

	return nil
}

// ImagePaths returns paths of all images of place p. They are
// specific to your machine!
func (p *Place) ImagePaths() []string {
	return p.imagePaths
}
