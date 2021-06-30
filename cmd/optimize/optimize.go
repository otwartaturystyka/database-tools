package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jdeng/goheif"
)

var noIcons bool

func init() {
	log.SetFlags(0)
	flag.BoolVar(&noIcons, "no-icons", false, "don't optimize icons")
}

func main() {
	flag.Parse()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("optimize: failed to get current working directory")
	}
	placeID := filepath.Base(currentDir)

	// Make srcPath - either .jpg or .heic
	originalIconPath := fmt.Sprintf("images/original/ic_%s.jpg", placeID)
	_, err = os.Stat(originalIconPath)
	if err != nil {
		if os.IsNotExist(err) {
			originalIconPath = fmt.Sprintf("images/original/ic_%s.heic", placeID)
		} else {
			log.Fatalf("optimize %s: failed to stat %s: %v\n:", placeID, originalIconPath, err)
		}
	}

	err = verifyValidDirectoryStructure(placeID, originalIconPath)
	if err != nil {
		log.Fatalf("optimize %s: no valid directory structure: %v\n", placeID, err)
	}

	if !noIcons {
		compressedIconPath := fmt.Sprintf("images/compressed/ic_%s.webp", placeID)
		err = makeIcon(originalIconPath, compressedIconPath)
		if err != nil {
			log.Fatalf("optimize %s: failed to create optimized icon: %v\n", placeID, err)
		}
		fmt.Printf("optimize %s: created optimized icon\n", placeID)
	}

	dirEntries, err := os.ReadDir("images/original")
	if err != nil {
		log.Fatalf("optimize %s: failed to read images/original/ dir: %v\n", placeID, err)
	}

	for _, dirEntry := range dirEntries {
		if strings.HasPrefix(dirEntry.Name(), ".") {
			continue
		}

		if strings.HasPrefix(dirEntry.Name(), "ic_") {
			continue
		}

		fullName := dirEntry.Name()
		name := strings.TrimSuffix(fullName, filepath.Ext(fullName))

		// Make srcPath - either .jpg or .heic
		srcPath := fmt.Sprintf("images/original/%s.jpg", name)
		_, err := os.Stat(srcPath)
		if err != nil {
			if os.IsNotExist(err) {
				srcPath = fmt.Sprintf("images/original/%s.heic", name)
			} else {
				log.Fatalf("optimize %s: failed to stat %s: %v\n:", placeID, srcPath, err)
			}
		}

		dstPath := fmt.Sprintf("images/compressed/%s.webp", name)

		err = makeImage(srcPath, dstPath)
		if err != nil {
			log.Fatalf("optimize %s: failed to create optimized image %s: %v\n", placeID, name, err)
		}
		fmt.Printf("optimize %s: created optimized image %s\n", placeID, name)
	}
}

func verifyValidDirectoryStructure(placeID string, originalIconPath string) error {
	// Does the images/ directory even exist?
	_, err := os.Stat("images")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("images/ dir does not exist")
		} else {
			return fmt.Errorf("stat images/ dir: %w", err)
		}
	}

	// Does images/original/ directory exist?
	_, err = os.Stat("images/original")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("images/original/ dir does not exist")
		} else {
			return fmt.Errorf("stat images/original/ dir: %w", err)
		}
	}

	if !noIcons {
		// Does images/original/ have a JPG icon?
		_, err = os.Stat(originalIconPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("icon at %s does not exist", originalIconPath)
			} else {
				return fmt.Errorf("stat icon at %s: %w", originalIconPath, err)
			}
		}

		// Is that JPG icon 1024x1024?
		// TODO
		w, h, err := getImageDimensions(originalIconPath)
		if err != nil {
			return fmt.Errorf("get image dimensions: %w", err)
		}

		if w != 1024 && h != 1024 {
			return fmt.Errorf("dimensions of %s are not 1024x1024", originalIconPath)
		}
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
			return fmt.Errorf("stat images/compressed/ dir: %w", err)
		}
	}

	return nil
}

/// MakeIcons creates a 512x512 WEBP version of a standard 1024x1024 JPG icon.
func makeIcon(srcPath string, dstPath string) error {
	cmd := exec.Command("magick", srcPath, "-resize", "512x512", dstPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("run ImageMagick for icon at %s: %w", srcPath, err)
	}

	return nil
}

// magick lesny_przystanek_1.heic -resize 50% -quality 75 lesny_przystanek_1.webp
func makeImage(srcPath string, dstPath string) error {
	cmd := exec.Command("magick", srcPath, "-resize", "25%", "-quality", "75", dstPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("run ImageMagick for image at %s: %w", srcPath, err)
	}

	return nil
}

func getImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}

	ext := filepath.Ext(imagePath)
	var config image.Config

	if ext == ".heic" {
		config, err = goheif.DecodeConfig(file)
		if err != nil {
			return 0, 0, err
		}
	} else {
		config, _, err = image.DecodeConfig(file)
		if err != nil {
			return 0, 0, err
		}
	}

	return config.Width, config.Height, nil
}
