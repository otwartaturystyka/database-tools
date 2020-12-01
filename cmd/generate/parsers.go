package main

import (
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
	var sections []internal.Section

	walker := func(path string, info os.FileInfo, err error) error {
		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}
		// Jump 2 levels down, to the track's directory.
		os.Chdir(path)

		var section internal.Section
		err = section.Parse(lang)
		if err != nil {
			return err
		}
		os.Chdir("../..")

		sections = append(sections, section)
		return nil
	}
	err := filepath.Walk("sections", walker)

	return sections, err
}

func parseTracks(lang string) ([]internal.Track, error) {
	var tracks []internal.Track

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
	var stories []internal.Story

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
	var dayrooms []internal.Dayroom

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
