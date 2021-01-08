package main

import (
	"os"
	"path/filepath"

	"golang.org/x/image/webp"

	"github.com/bartekpacia/database-tools/internal"
	"github.com/bbrks/go-blurhash"
	"github.com/pkg/errors"
)

// ParseMeta parses metadata for the generated datafile of ID regionID.
func parseMeta(regionID string, lang string) (*internal.Meta, error) {
	datafilePath := filepath.Join("generated", regionID)

	if err := os.Chdir(datafilePath); err != nil {
		return nil, errors.Wrapf(err, "failed to chdir into generated datafile's at %s", datafilePath)
	}

	var meta internal.Meta
	err := meta.ParseFromGenerated()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse meta from a generated datafile's data.json at %s", datafilePath)
	}

	err = os.Chdir("../../")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to exit (aka go 2 dirs up) generated datafile's path at %s ", datafilePath)
	}

	return &meta, nil
}

func makeThumbBlurhash(regionID string) (blur string, err error) {
	file, err := os.Open("database/" + regionID + "/meta/thumb_mini.webp")
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
