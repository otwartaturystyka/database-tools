package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/internal"
)

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
	if err != nil {
		log.Fatalln("generate:", err)
	}

	return tracks, nil
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
	if err != nil {
		log.Fatalln("generate:", err)
	}

	return stories, nil
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
	if err != nil {
		log.Fatalln("generate:", err)
	}

	return dayrooms, nil
}
