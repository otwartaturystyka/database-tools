package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	regionID      string
	justIcons     bool
	splitPaths    bool
	justFilenames bool
)

type entry struct {
}

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region id")
	flag.BoolVar(&justIcons, "just-icons", false, "whether to consider only icons (files starting with ic_)")
	flag.BoolVar(&splitPaths, "split-paths", false, "whether to split filepaths")
	flag.BoolVar(&justFilenames, "just-filenames", false, "whether to print only filenames")
}

func main() {
	flag.Parse()

	if regionID == "" {
		log.Fatalln("walk: regionID is empty")
	}

	jpegs := 0
	pngs := 0
	webps := 0
	walker := func(path string, info os.FileInfo, err error) error {
		level := strings.Count(path, "/")
		if level != 6 || strings.Contains(path, "/.git/") {
			return nil
		}
		ext := filepath.Ext(path)

		if justIcons && !strings.Contains(path, "/ic_") {
			return nil
		}

		if ext == ".jpg" || ext == ".jpeg" {
			jpegs++
		} else if ext == ".png" {
			pngs++
		} else if ext == ".webp" {
			webps++
		} else {
			return nil
		}

		sizeMB := float32(info.Size()) / 1000 / 1000

		// if sizeMB < 2 {
		// 	return nil
		// }

		splitties := strings.Split(path, "/")
		filename := splitties[len(splitties)-1]
		path = strings.TrimSuffix(path, filename) + " " + filename

		if splitPaths {
			fmt.Printf("walk: %.2f MB %s\n", sizeMB, path)
		}

		if justFilenames {
			fmt.Printf("walk: %.2f MB %s\n", sizeMB, filename)
		}

		return nil
	}

	filepath.Walk("database/"+regionID, walker)

	fmt.Printf("walk: %d jpegs, %d pngs, %d webps \n", jpegs, pngs, webps)
}
