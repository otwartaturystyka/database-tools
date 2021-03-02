package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/readers"
	"github.com/pkg/errors"
)

type Action struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

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
	WebsiteURL  string   `json:"website_url"`
	FacebookURL string   `json:"facebook_url"`
	Headers     []string `json:"headers"`
	Content     []string `json:"content"`
	Actions     []Action `json:"actions"`
	Images      []string `json:"images"`
	imagePaths  []string
}

// Parse parses place data from its directory and assigns
// it to track pointed to by p. It must be used directly
// in the place's directory.
func (p *Place) Parse(lang string) error {
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
		return errors.Wrapf(err, "make image paths for place %s", p.ID)
	}

	// Content
	name, err := readers.ReadFromFile("content/" + lang + "/name.txt")
	if err != nil {
		return err
	}
	p.Name = strings.TrimSuffix(string(name), "\n")

	quickInfo, err := readers.ReadFromFile("content/" + lang + "/quick_info.txt")
	if err != nil {
		return err
	}
	p.QuickInfo = strings.TrimSuffix(string(quickInfo), "\n")

	overview, err := readers.ReadFromFile("content/" + lang + "/overview.txt")
	if err != nil {
		return err
	}
	p.Overview = strings.TrimSuffix(string(overview), "\n")

	// Headers and content
	p.Headers = make([]string, 0)
	p.Content = make([]string, 0)
	p.Actions = make([]Action, 0)
	for i := 0; true; i++ {
		textFilePath := filepath.Join("content", lang, fmt.Sprintf("text_%d.txt", i))
		textFile, err := os.Open(textFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("file %s of place %s does not exist (most probably, this place does not have additional content)\n", textFilePath, p.ID)
				break
			}

			fmt.Printf("failed to open file %s: %v\n", textFilePath, err)
		}

		header, content, err := readers.ReadTextualData(textFile)
		if err != nil {
			return err
		}
		err = textFile.Close()
		if err != nil {
			return err
		}

		p.Headers = append(p.Headers, strings.TrimSuffix(header, "\n"))
		p.Content = append(p.Content, strings.TrimSuffix(content, "\n"))

		err = p.makeActions(lang)
		if err != nil {
			return errors.Wrapf(err, "make actions for place %s", p.ID)
		}
	}

	return nil
}

func (p *Place) makeActions(lang string) error {
	actionValuesFile, err := readers.ReadFromFile("actions.json")
	if err != nil {
		fmt.Printf("file %s of place %s does not exist (most probably, this place does not have any actions)\n", "actions.json", p.ID)
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

			return errors.WithStack(err)
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
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "make image paths: failed to get working dir")
	}

	var qualityDir string
	if quality == Compressed {
		qualityDir = "compressed"
	} else if quality == Original {
		qualityDir = "original"
	}

	// s.Images were set when the story was parsed from its JSON.
	for _, image := range p.Images {
		absImagePath := filepath.Join(wd, "images", qualityDir, image+".webp")

		if _, err := os.Stat(absImagePath); err != nil {
			return errors.Errorf("image at %s does not exist!\n", absImagePath)
		}

		p.imagePaths = append(p.imagePaths, absImagePath)
	}

	// Add icon.
	absIconPath := filepath.Join(wd, "images", qualityDir, p.Icon+".webp")
	p.imagePaths = append(p.imagePaths, absIconPath)

	return nil
}

// ImagePaths returns paths of all images of place p. They are
// specific to your machine!
func (p *Place) ImagePaths() []string {
	return p.imagePaths
}
