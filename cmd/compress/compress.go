package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("compress: ")
}

func Compress(regionID string, verbose bool) {
	_, err := os.Stat("compressed/")
	if os.IsNotExist(err) {
		err = os.Mkdir("compressed", 0755)
		if err != nil {
			log.Fatalf("compress: error creating compressed directory: %v\n", err)
		}
	}

	zipFile, err := os.Create(filepath.Join("compressed", regionID+".zip"))
	if err != nil {
		log.Fatalln("compress: failed to create zip file")
	}
	defer zipFile.Close()

	sourceDatafilePath := filepath.Join("generated", regionID)

	info, err := os.Stat(sourceDatafilePath)
	if os.IsNotExist(err) {
		log.Fatalf("compress: datafile in \"generated\" directory for %s doesn't exist", regionID)
	}

	if !info.IsDir() {
		log.Fatalf("compress: error: datafile %s is not a directory\n", regionID)
	}

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	i := 0
	walker := func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.Name() == "." {
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

		// Crop the "generated/" part to remove it from the result .zip.
		components := strings.SplitAfterN(path, "/", 2)
		writer, err := zipWriter.Create(components[1])
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

	err = filepath.Walk(sourceDatafilePath, walker)
	if err != nil {
		log.Fatalf("compress: error while walking %s: %v\n", regionID, err)
	}

	fmt.Println("compress: successfully compressed datafile", regionID)
}
