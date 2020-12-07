package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/image/webp"
)

var (
	regionID      string
	minSize       float64
	sortBy        string
	icons         bool
	splitPaths    bool
	justFilenames bool
	count         bool
)

type entry struct {
	filename string
	path     string
	sizeMB   float64
	width    int
	height   int
}

func (e *entry) dimens() string {
	return fmt.Sprintf("%dx%d", e.width, e.height)
}

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.StringVar(&regionID, "region-id", "", "region id")
	flag.Float64Var(&minSize, "min-size", 2, "min size of images to list")
	flag.StringVar(&sortBy, "sort-by", "count", "sort by \"count\" or \"ratio\"")
	flag.BoolVar(&icons, "icons", false, "whether to list icons (files starting with ic_)")
	flag.BoolVar(&splitPaths, "split-paths", false, "whether to split filepaths")
	flag.BoolVar(&justFilenames, "just-filenames", false, "whether to print only filenames")
	flag.BoolVar(&count, "count", false, "whether to show number next to each entry")
}

func main() {
	flag.Parse()

	if regionID == "" {
		log.Fatalln("walk: regionID is empty")
	}

	jpegs := 0
	pngs := 0
	webps := 0
	entries, stats := discover(&jpegs, &pngs, &webps)

	fmt.Println(len(entries))

	// Sort by age, keeping original order or equal elements.
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].sizeMB > entries[j].sizeMB
	})

	for _, e := range entries {
		if splitPaths {
			fmt.Printf("walk: %.2f MB %s\n", e.sizeMB, e.path)
		}

		if justFilenames {
			fmt.Printf("walk: %.2f MB %s %s\n", e.sizeMB, e.dimens(), e.filename)
		}

		if !splitPaths && !justFilenames {
			fmt.Printf("walk: %.2f MB %s %s\n", e.sizeMB, e.dimens(), e.path+e.filename)
		}
	}

	type pair struct {
		dimens string
		count  int
		ratio  float32
	}

	var pairs []pair
	for dimens, count := range stats {
		s := strings.Split(dimens, "x")
		w, _ := strconv.Atoi(s[0])
		h, _ := strconv.Atoi(s[1])

		ratio := float32(h) / float32(w)
		pair := pair{dimens: dimens, count: count, ratio: ratio}
		pairs = append(pairs, pair)
	}

	if sortBy == "count" {
		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].count > pairs[j].count
		})
	} else if sortBy == "ratio" {
		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].ratio > pairs[j].ratio
		})
	}

	for _, pair := range pairs {
		fmt.Printf("walk: %d images of size %s and aspect ratio %.2f\n", pair.count, pair.dimens, pair.ratio)
	}

	total := jpegs + pngs + webps
	fmt.Printf("walk: %d jpegs, %d pngs, %d webps (%d total) \n", jpegs, pngs, webps, total)
}

func discover(jpegs *int, pngs *int, webps *int) (entries []entry, stats map[string]int) {
	stats = make(map[string]int, 0)

	walker := func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "/.git/") {
			return nil
		}

		if !strings.Contains(path, "/original/") && !strings.Contains(path, "/compressed/") {
			return nil
		}
		ext := filepath.Ext(path)

		if !icons && strings.Contains(path, "/ic_") {
			return nil
		}

		w, h := 0, 0
		var file *os.File
		if ext == ".jpg" || ext == ".jpeg" {
			*(jpegs)++
			file, err = os.Open(path)
			if err != nil {
				log.Fatalln("walk: error opening file:", err)
			}

			cfg, err := jpeg.DecodeConfig(file)
			if err != nil {
				log.Fatalln("walk: error decoding image config:", path, err)
			}

			w = cfg.Width
			h = cfg.Height
		} else if ext == ".png" {
			*(pngs)++
			file, err = os.Open(path)
			if err != nil {
				log.Fatalln("walk: error opening file:", err)
			}

			cfg, err := png.DecodeConfig(file)
			if err != nil {
				log.Fatalln("walk: error decoding image config:", err)
			}

			w = cfg.Width
			h = cfg.Height
		} else if ext == ".webp" {
			*(webps)++
			file, err = os.Open(path)
			if err != nil {
				log.Fatalln("walk: error opening file:", err)
			}

			cfg, err := webp.DecodeConfig(file)
			if err != nil {
				log.Fatalln("walk: error decoding image config:", err)
			}

			w = cfg.Width
			h = cfg.Height
		} else {
			fmt.Println("WTFFFFFFF", path)
			return nil
		}
		defer file.Close()

		sizeMB := float64(info.Size()) / 1000 / 1000

		if sizeMB < minSize {
			return nil
		}

		splitties := strings.Split(path, "/")
		filename := splitties[len(splitties)-1]
		justPath := strings.TrimSuffix(path, filename)

		if h > w {
			fmt.Println(path, filename)
		}

		entry := entry{path: justPath, filename: filename, sizeMB: sizeMB, width: w, height: h}
		entries = append(entries, entry)
		stats[entry.dimens()]++
		return nil
	}

	filepath.Walk("database/"+regionID, walker)

	return
}
