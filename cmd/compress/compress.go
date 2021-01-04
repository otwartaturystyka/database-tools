package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	regionID string
	verbose  bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&regionID, "region-id", "", "region which datafile should be compressed")
	flag.BoolVar(&verbose, "verbose", false, "true for extensive logging")
	flag.Parse()

	if regionID == "" {
		log.Fatalln("compress: error: regionID is empty")
	}
}

func main() {
	_, err := os.Stat("generated/")
	if os.IsNotExist(err) {
		err = os.Mkdir("generated", 0755)
		if err != nil {
			log.Fatalf("compress: error creating generated directory: %v\n", err)
		}
	}

	zipFile, err := os.Create("compressed/" + regionID + ".zip")
	defer zipFile.Close()

	os.Chdir("database/")
	wd, _ := os.Getwd()
	if verbose {
		fmt.Println("compress: changed working directory to", wd)
	}

	_, err = os.Stat(regionID)
	if os.IsNotExist(err) {
		log.Fatalf("compress: datafile directory for %s doesn't exist", regionID)
	}

	info, err := os.Stat(regionID)
	if os.IsNotExist(err) {
		log.Fatalf("compress: error: directory for region %s doesn't exist\n", regionID)
	}

	if !info.IsDir() {
		log.Fatalf("compress: error: datafile %s is not a directory\n", regionID)
	}

	w := zip.NewWriter(zipFile)
	defer w.Close()

	i := 0
	walker := func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.Name() == ".DS_Store" {
			return nil
		}

		if fileInfo.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := w.Create(path)
		if err != nil {
			log.Fatalf("compress: error creating a file in zip archive: %v\n", err)
		}

		if verbose {
			fmt.Printf("compress: compressing file %d at %s\n", i, path)
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			log.Fatalf("compress: error copying file: %v\n", err)
		}

		i++
		return nil
	}

	err = filepath.Walk(regionID, walker)
	if err != nil {
		log.Fatalf("compress: error while walking %s: %v\n", regionID, err)
	}

	os.Chdir("..")
	if verbose {
		fmt.Println("compress: changed working directory back")
	}

	fmt.Println("compress: successfully compressed datafile", regionID)
}
