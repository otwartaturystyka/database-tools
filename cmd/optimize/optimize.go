package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	log.SetFlags(0)
}

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("optimize: failed to get current working directory")
	}

	placeID := filepath.Base(currentDir)
	originalIconPath := fmt.Sprintf("images/original/ic_%s.jpg", placeID)
	compressedIconPath := fmt.Sprintf("images/compressed/ic_%s.webp", placeID)

	err = verifyValidDirectoryStructure(placeID, originalIconPath)
	if err != nil {
		log.Fatalf("optimize %s: no valid directory structure: %v\n", placeID, err)
	}

	_, err = os.Stat(compressedIconPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("optimize %s: compressed icon does not exist. It will be created\n", placeID)
			err = makeIcon(originalIconPath, compressedIconPath)
			if err != nil {
				log.Fatalf("optimize %s: failed to create images/compressed/ dir\n", placeID)
			}
			fmt.Printf("optimize %s: created compressed icon\n", placeID)
		} else {
			log.Fatalf("optimize %s: failed to stat images/ dir: %v\n", err, placeID)
		}
	}

	fileNames, err := os.ReadDir("images/original")
	if err != nil {
		log.Fatalf("optimize %s: failed to read images/original/ dir: %v\n", placeID, err)
	}

	for _, fileName := range fileNames {
		if strings.HasPrefix(fileName.Name(), "ic_") {
			continue
		}

		fmt.Println(fileName.Name())
	}
}

func verifyValidDirectoryStructure(placeID string, originalIconPath string) error {
	// Does the images/ directory even exist?
	_, err := os.Stat("images")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("images/ dir does not exist")
		} else {
			return errors.Wrap(err, "stat images")
		}
	}

	// Does images/original/ directory exist?
	_, err = os.Stat("images/original")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("images images/original/ dir does not exist")
		} else {
			return errors.Wrap(err, "stat images/original/")
		}
	}

	// Does images/original/ have a JPG icon?
	_, err = os.Stat(originalIconPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Wrapf(err, "%s does not exist", originalIconPath)
		} else {
			return errors.Wrapf(err, "stat %s", originalIconPath)
		}
	}

	// Is that JPG icon 1024x1024?
	// TODO
	w, h, err := getImageDimensions(originalIconPath)
	if err != nil {
		return errors.Wrap(err, "get image dimensions")
	}

	if w != 1024 && h != 1024 {
		return errors.Errorf("dimensions of %s are not 1024x1024", originalIconPath)
	}

	// Does images/compressed/ directory exist?
	_, err = os.Stat("images/compressed")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("images/compressed", 0755)
			if err != nil {
				return errors.New("create images/compressed/ dir")
			}
			fmt.Printf("optimize %s: images/compressed/ dir was created\n", placeID)
		} else {
			return errors.Wrap(err, "stat images/compressed/")
		}
	}

	return nil
}

/// MakeIcons creates a 512x512 WEBP version of a standard 1024x1024 JPG icon.
func makeIcon(srcPath string, dstPath string) error {
	err := exec.Command("magick", srcPath, "-resize", "512x512", dstPath).Run()
	if err != nil {
		return errors.Wrapf(err, "optimize icon %s", srcPath)
	}

	return nil
}

// magick lesny_przystanek_1.heic -resize 50% -quality 75 lesny_przystanek_1.webp
func makeImage(srcPath string, dstPath string) error {
	err := exec.Command("magick", srcPath, "-resize 25%", "-quality 75", dstPath).Run()
	if err != nil {
		return errors.Wrapf(err, "optimize icon %s", srcPath)
	}

	return nil
}

func getImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return image.Width, image.Height, nil
}
