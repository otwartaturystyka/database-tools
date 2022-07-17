package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/opentouristics/database-tools/models"
)

func getCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %v", err)
	}

	hash := strings.TrimSpace(string(out))
	return hash, nil
}

func getCommitTag() (string, error) {
	cmd := exec.Command("git", "tag", "--points-at")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git tag: %v", err)
	}

	tag := strings.TrimSpace(string(out))
	return tag, nil
}

func parseMeta(lang string) (meta models.Meta, err error) {
	os.Chdir("meta")

	err = meta.Parse(lang)
	if err != nil {
		err = fmt.Errorf("parse meta: %w", err)
		return
	}

	os.Chdir("..")

	commitHash, err := getCommitHash()
	if err != nil {
		err = fmt.Errorf("get commit hash: %v", err)
		return
	}
	meta.CommitHash = commitHash

	commitTag, err := getCommitTag()
	if err != nil {
		err = fmt.Errorf("get commit tag: %v", err)
		return
	}
	if commitTag == "" {
		meta.CommitTag = nil
	} else {
		meta.CommitTag = &commitTag
	}

	return
}

func parseSections(lang string, verbose bool) ([]models.Section, error) {
	sections := make([]models.Section, 0)

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("start walk section %s: %w", path, err)
		}

		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}
		// Jump 2 levels down, to the sections's directory.
		os.Chdir(path)

		var section models.Section
		err = section.Parse(lang, verbose)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		os.Chdir("../..")

		sections = append(sections, section)
		return nil
	}
	err := filepath.Walk("sections", walker)
	if err != nil {
		return sections, fmt.Errorf("walk sections: %w", err)
	}

	return sections, err
}

func parseTracks(lang string) ([]models.Track, error) {
	tracks := make([]models.Track, 0)

	if _, err := os.Stat("tracks"); os.IsNotExist(err) {
		return tracks, nil
	}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("parse track %s: %w", path, err)
		}

		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}
		// Jump 2 levels down, to the track's directory.
		os.Chdir(path)

		var track models.Track
		err = track.Parse(lang)

		tracks = append(tracks, track)
		os.Chdir("../..")

		return err
	}
	err := filepath.Walk("tracks", walker)

	return tracks, err
}

func parseStories(lang string) ([]models.Story, error) {
	stories := make([]models.Story, 0)

	if _, err := os.Stat("stories"); os.IsNotExist(err) {
		return stories, nil
	}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("parse story %s: %w", path, err)
		}

		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}
		// Jump 2 levels down, to the story's directory.
		os.Chdir(path)

		var story models.Story
		err = story.Parse(lang)

		stories = append(stories, story)
		os.Chdir("../..")

		return err
	}

	err := filepath.Walk("stories", walker)

	return stories, err
}
