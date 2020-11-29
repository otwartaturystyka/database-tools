package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bartekpacia/database-tools/internal"
)

var (
	regionID string
	language string
	quality  int
	verbose  bool
)

func check(err error) {
	if err != nil {
		log.Fatalln("generate:", err)
	}
}

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.StringVar(&language, "lang", "pl", "language of text in the datafile")
	flag.IntVar(&quality, "quality", 1, "quality of photos in the datafile")
}

func main() {
	flag.Parse()

	var datafile internal.Datafile

	if regionID == "" {
		log.Fatalln("generate: regionID is empty")
	}

	if quality != 1 && quality != 2 {
		log.Fatalln("generate: quality is not 1 or 2")
	}

	err := os.Chdir("database/" + regionID)
	if err != nil {
		log.Fatalln("generate:", err)
	}

	tracks, err := parseTracks(language)
	check(err)
	datafile.Tracks = tracks

	dayrooms, err := parseDayrooms(language)
	check(err)
	datafile.Dayrooms = dayrooms

	os.Chdir("../..")

	err = os.RemoveAll("generated/" + regionID)
	check(err)

	err = os.Mkdir("generated/"+regionID+"/", 0755)
	check(err)

	dataJSONFile, err := os.Create("generated/" + regionID + "/data.json")
	check(err)

	b, err := json.MarshalIndent(datafile, "", "	")
	check(err)

	n, err := dataJSONFile.Write(b)
	check(err)

	fmt.Printf("generate: wrote %d KB to data.json file\n", n/1024)
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
	if err != nil {
		log.Fatalln("generate:", err)
	}

	return tracks, nil
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
