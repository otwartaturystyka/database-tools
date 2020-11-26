package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	regionID string
	verbose  bool
)

func init() {
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be compressed")
	flag.BoolVar(&verbose, "verbose", false, "true for extensive logging")
}

func main() {
	flag.Parse()

	fmt.Println(regionID)

	if regionID == "" {
		log.Fatalln("compress: error: regionID is empty")
	}

	datafilePath := "database/" + regionID
	info, err := os.Stat(datafilePath)
	if err != nil {
		log.Fatalf("compress: error: directory for %s not found: %v\n", regionID, err)
	}

	if !info.IsDir() {
		log.Fatalf("compress: error: datafile %s is not a directory\n", regionID)
	}

	outputFile, err := os.Create(regionID + ".zip")
	if err != nil {
		log.Fatalf("compress: error: failed to create a zip output file %s: %v\n", regionID, err)
	}
	defer outputFile.Close()

	w := zip.NewWriter(outputFile)
	defer w.Close()

	i := 0
	walker := func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.Name() == ".DS_Store" {
			fmt.Println("compress: encountered DS_Store")
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Println(fmt.Sprint(i), path)

		i++

		return nil
	}

	err = filepath.Walk(datafilePath, walker)
	if err != nil {
		log.Fatalf("compress: error while walking %s: %v\n", datafilePath, err)
	}

}
