// Package optimize implements a simple image optimization functionality.
package optimize

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jdeng/goheif"
)

// Optimize creates optimized versions of images from images in the place's "original" directory.
// placePath must point to a valid place.
func Optimize(placeID string, noIcons bool, verbose bool) error {
	// Make srcPath - either .jpg or .heic
	originalIconPath := fmt.Sprintf("images/original/ic_%s.jpg", placeID)
	_, err := os.Stat(originalIconPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			originalIconPath = fmt.Sprintf("images/original/ic_%s.heic", placeID)
		} else {
			return fmt.Errorf("stat %s: %v", originalIconPath, err)
		}
	}

	err = verifyValidDirectoryStructure(placeID, originalIconPath, noIcons)
	if err != nil {
		return fmt.Errorf("no valid directory structure: %v", err)
	}

	if !noIcons {
		compressedIconPath := fmt.Sprintf("images/compressed/ic_%s.webp", placeID)
		err = makeIcon(originalIconPath, compressedIconPath)
		if err != nil {
			return fmt.Errorf("make optimized icon: %v", err)
		}

		if verbose {
			fmt.Printf("optimize %s: created optimized icon\n", placeID)
		}
	}

	dirEntries, err := os.ReadDir("images/original")
	if err != nil {
		return fmt.Errorf("read images/original/ directory: %v", err)
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
			if errors.Is(err, os.ErrNotExist) {
				srcPath = fmt.Sprintf("images/original/%s.heic", name)
			} else {
				return fmt.Errorf("stat %s: %v", srcPath, err)
			}
		}

		dstPath := fmt.Sprintf("images/compressed/%s.webp", name)

		err = makeImage(srcPath, dstPath)
		if err != nil {
			return fmt.Errorf("create optimized image %s: %v", name, err)
		}

		if verbose {
			fmt.Printf("created optimized image %s\n", name)
		}
	}

	return nil
}

func verifyValidDirectoryStructure(placeID string, originalIconPath string, noIcons bool) error {
	// Check if images/ directory exists
	_, err := os.Stat("images")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("images/ dir does not exist")
		} else {
			return fmt.Errorf("stat images/ dir: %w", err)
		}
	}

	// Check if images/original/ directory exists
	_, err = os.Stat("images/original")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("images/original/ dir does not exist")
		} else {
			return fmt.Errorf("stat images/original/ dir: %w", err)
		}
	}

	if !noIcons {
		// Check if images/original/ directory has a JPG icon.
		_, err = os.Stat(originalIconPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("icon at %s does not exist", originalIconPath)
			} else {
				return fmt.Errorf("stat icon at %s: %w", originalIconPath, err)
			}
		}

		// TODO: check if that JPG icon is 1024x1024
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
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir("images/compressed", 0o755)
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

// MakeIcons creates a compressed 512x512 WEBP version of an original 1024x1024 JPG or HEIC icon.
func makeIcon(srcPath string, dstPath string) error {
	cmd := exec.Command("magick", srcPath, "-resize", "512x512", dstPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("run ImageMagick for icon at %s: %w", srcPath, err)
	}

	return nil
}

// MakeImage creates a compressed WEBP version of an original JPEG or HEIC image.
// The compressed image has 4 times smaller resolution and also has decreased quality.
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
