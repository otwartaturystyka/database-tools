package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bartekpacia/database-tools/internal"
	"github.com/bbrks/go-blurhash"
	"golang.org/x/image/webp"
)

// ParseMeta parses metadata for the generated datafile of ID regionID.
func parseMeta(regionID string, lang string) (*internal.Meta, error) {
	datafilePath := filepath.Join("generated", regionID)

	if err := os.Chdir(datafilePath); err != nil {
		return nil, fmt.Errorf("chdir into generated datafile's dir at %s: %w", datafilePath, err)
	}

	var meta internal.Meta
	err := meta.ParseFromGenerated()
	if err != nil {
		return nil, fmt.Errorf("parse meta from generated datafile's data.json at %s: %w", datafilePath, err)
	}

	err = os.Chdir("../../")
	if err != nil {
		return nil, fmt.Errorf("exit (=chdir 2 dirs up) generated datafile's dir at %s: %w", datafilePath, err)
	}

	return &meta, nil
}

func makeThumbBlurhash(regionID string) (blur string, err error) {
	file, err := os.Open(filepath.Join("database", regionID, "meta", "thumb_mini.webp"))
	if err != nil {
		return "", err
	}

	thumbImage, err := webp.Decode(file)
	if err != nil {
		return "", err
	}

	blur, err = blurhash.Encode(4, 3, thumbImage)
	if err != nil {
		return "", nil
	}

	return
}
