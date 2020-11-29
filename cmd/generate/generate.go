package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be uploaded")
	flag.StringVar(&language, "lang", "pl", "language of text in the datafile")
	flag.IntVar(&quality, "quality", 1, "quality of photos in the datafile")
}

func main() {
	flag.Parse()

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

	dayrooms, err := parseDayrooms(language)
	if err != nil {
		log.Fatalln("generate:", err)
	}

	for _, dayroom := range dayrooms {
		fmt.Println(dayroom.Name, dayroom.Lat, dayroom.Lng)
	}
}

func parseDayrooms(language string) (dayrooms []internal.Dayroom, err error) {
	walker := func(path string, info os.FileInfo, err error) error {
		level := strings.Count(path, "/")
		if level != 1 {
			return nil
		}

		var dayroom internal.Dayroom

		os.Chdir(path)

		jsonFile, err := os.Open("data.json")
		if err != nil {
			log.Fatalln("generate:", err)
		}

		nameFile, err := os.Open("content/" + language + "/name.txt")
		if err != nil {
			log.Fatalln("generate:", err)
		}

		name, err := ioutil.ReadAll(nameFile)
		if err != nil {
			log.Fatalln("generate:", err)
		}
		dayroom.Name = string(name)

		data, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalln("generate:", err)
		}

		err = json.Unmarshal(data, &dayroom)
		if err != nil {
			log.Fatalln("generate:", err)
		}

		dayrooms = append(dayrooms, dayroom)

		os.Chdir("../..")
		return nil
	}

	err = filepath.Walk("dayrooms", walker)
	if err != nil {
		log.Fatalln("generate:", err)
	}

	return dayrooms, nil
}
