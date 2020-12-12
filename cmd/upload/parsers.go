package main

import (
	"os"
	"path/filepath"

	"golang.org/x/image/webp"

	"github.com/bartekpacia/database-tools/internal"
	"github.com/bbrks/go-blurhash"
	"github.com/pkg/errors"
)

func parseMeta(regionID string, lang string) (*internal.Meta, error) {
	metaPath := filepath.Join("database", regionID, "meta")

	if err := os.Chdir(metaPath); err != nil {
		return nil, errors.Wrapf(err, "failed to chdir into metaPath at %s", metaPath)
	}

	var meta internal.Meta
	meta.Parse(lang)

	if err := os.Chdir("../../../"); err != nil {
		return nil, errors.Wrapf(err, "failed to exit (aka go 3 dirs up) metaPath at %s ", metaPath)
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
