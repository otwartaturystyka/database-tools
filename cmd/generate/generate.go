package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	dayroomsInfo, err := os.Stat("dayrooms")
	if err != nil {
		log.Fatalln("generate:", err)
	}

	fmt.Println(dayroomsInfo)

	fmt.Println("generate: not implemented yet")
}
