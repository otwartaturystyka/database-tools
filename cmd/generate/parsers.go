package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/internal"
	"github.com/pkg/errors"
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
		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}
		// Jump 2 levels down, to the sections's directory.
		os.Chdir(path)

		var section internal.Section
		err = section.Parse(lang)
		if err != nil {
			return errors.Wrapf(err, "parse section \"%s\"", path)
		}
		os.Chdir("../..")

		sections = append(sections, section)
		return nil
	}
	err := filepath.Walk("sections", walker)
	if err != nil {
		return sections, errors.Wrapf(err, "walk sections")
	}

	return sections, err
}

func parseTracks(lang string) ([]internal.Track, error) {
	tracks := make([]internal.Track, 0)

	walker := func(path string, info os.FileInfo, err error) error {
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

	walker := func(path string, info os.FileInfo, err error) error {
		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}
		// Jump 2 levels down, to dayrooms directory.
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

func parseDayrooms(lang string) ([]internal.Dayroom, error) {
	dayrooms := make([]internal.Dayroom, 0)

	walker := func(path string, info os.FileInfo, err error) error {
		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}
		// Jump 2 levels down, to dayrooms directory.
		os.Chdir(path)

		var dayroom internal.Dayroom
		err = dayroom.Parse(lang)

		dayrooms = append(dayrooms, dayroom)
		os.Chdir("../..")

		return err
	}

	err := filepath.Walk("dayrooms", walker)

	return dayrooms, err
}
