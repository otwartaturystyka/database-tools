// Package compress handles compressing the generated datafile directory.
package compress

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Compress takes a generated directory of region's datafile and creates a zip archive out of it.
func Compress(regionID string, verbose bool) error {
	_, err := os.Stat("compressed/")
	if os.IsNotExist(err) {
		err = os.Mkdir("compressed", 0755)
		if err != nil {
			return fmt.Errorf("failed to create compressed directory: %v", err)
		}
	}

	zipFilePath := filepath.Join("compressed", regionID+".zip")
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	sourceDatafilePath := filepath.Join("generated", regionID)

	info, err := os.Stat(sourceDatafilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("compress: datafile in \"generated\" directory for %s doesn't exist", regionID)
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", sourceDatafilePath)
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
			return fmt.Errorf("create a file in zip archive: %v", err)
		}

		if verbose {
			fmt.Printf("compressing file %d at %s\n", i, path)
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("copy %s: %v", path, err)
		}

		i++
		return nil
	}

	err = filepath.Walk(sourceDatafilePath, walker)
	if err != nil {
		return fmt.Errorf("walk %s: %v", sourceDatafilePath, err)
	}

	fmt.Println("successfully compressed datafile", regionID)

	return nil
}
