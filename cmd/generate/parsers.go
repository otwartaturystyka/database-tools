package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/internal"
)

func parseMeta(lang string) (internal.Meta, error) {
	os.Chdir("meta")

	var meta internal.Meta
	err := meta.Parse(lang)

	os.Chdir("..")

	return meta, err
}

func parseSections(lang string) ([]internal.Section, error) {
	sections := make([]internal.Section, 0)

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

		var section internal.Section
		err = section.Parse(lang)
		if err != nil {
			return fmt.Errorf("parse section %s: %w", path, err)
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

func parseTracks(lang string) ([]internal.Track, error) {
	tracks := make([]internal.Track, 0)

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

		var track internal.Track
		err = track.Parse(lang)

		tracks = append(tracks, track)
		os.Chdir("../..")

		return err
	}
	err := filepath.Walk("tracks", walker)

	return tracks, err
}

func parseStories(lang string) ([]internal.Story, error) {
	stories := make([]internal.Story, 0)

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

		var story internal.Story
		err = story.Parse(lang)

		stories = append(stories, story)
		os.Chdir("../..")

		return err
	}

	err := filepath.Walk("stories", walker)

	return stories, err
}
