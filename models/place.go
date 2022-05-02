package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/formatters"
	"github.com/bartekpacia/database-tools/readers"
)

// Place represents single place in real world.
type Place struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Section     string   `json:"section"`
	Icon        string   `json:"icon"`
	QuickInfo   string   `json:"quick_info"`
	Overview    string   `json:"overview"`
	Lat         float32  `json:"lat"`
	Lng         float32  `json:"lng"`
	WebsiteURL  *string  `json:"website_url"`
	FacebookURL *string  `json:"facebook_url"`
	Headers     []string `json:"headers"`
	Content     []string `json:"content"`
	Actions     []Action `json:"actions"`
	Images      []string `json:"images"`
	imagePaths  []string
}

// Parse parses place data from its directory and assigns
// it to track pointed to by p. It must be used directly
// in the place's directory.
func (p *Place) Parse(lang string, verbose bool) error {
	// Technicla metadata
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
	name, err := readers.ReadFromFile(filepath.Join("content", lang, "name.txt"))
	if err != nil {
		return err
	}
	p.Name = formatters.ToContent(string(name))

	quickInfo, err := readers.ReadFromFile(filepath.Join("content", lang, "quick_info.txt"))
	if err != nil {
		return err
	}
	p.QuickInfo = formatters.ToContent(string(quickInfo))

	overview, err := readers.ReadFromFile(filepath.Join("content", lang, "overview.txt"))
	if err != nil {
		return err
	}
	p.Overview = formatters.ToContent(string(overview))

	// Headers and content
	p.Headers = make([]string, 0)
	p.Content = make([]string, 0)
	for i := 0; true; i++ {
		textFilePath := filepath.Join("content", lang, fmt.Sprintf("text_%d.txt", i))
		textFile, err := os.Open(textFilePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if verbose {
					fmt.Printf(
						"file %s of place %s does not exist (most probably, "+
							"this place does not have any content, so the file "+
							"does not exist)\n",
						textFilePath,
						p.ID,
					)
				}
				break
			}

			fmt.Printf("failed to open file %s: %v\n", textFilePath, err)
		}
		defer textFile.Close()

		sectionText, err := io.ReadAll(textFile)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", textFilePath, err)
		}

		header, content := formatters.ToSection(string(sectionText))

		p.Headers = append(p.Headers, header)
		p.Content = append(p.Content, content)
	}

	// Actions
	p.Actions = make([]Action, 0)
	err = p.makeActions(lang, verbose)
	if err != nil {
		return fmt.Errorf("make actions for place %s: %w", p.ID, err)
	}

	return nil
}

func (p *Place) makeActions(lang string, verbose bool) error {
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

	// Read action names from a valid language file
	actionNames := make([]string, 0)
	for i := 0; i < len(actionValues); i++ {
		actionFilePath := filepath.Join("content", lang, fmt.Sprintf("action_%d.txt", i))
		b, err := readers.ReadFromFile(actionFilePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("file %s of place %s does not exist (most probably, this means the translation is missing)\n", actionFilePath, p.ID)
				break
			}

			return err
		}

		actionName := strings.TrimSuffix(string(b), "\n")
		actionNames = append(actionNames, actionName)
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
