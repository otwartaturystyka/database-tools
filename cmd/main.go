package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/opentouristics/database-tools/cmd/compress"
	"github.com/opentouristics/database-tools/cmd/generate"
	"github.com/opentouristics/database-tools/cmd/optimize"
	"github.com/opentouristics/database-tools/cmd/upload"
	"github.com/opentouristics/database-tools/models"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFlags(0)
}

var generateCommand = cli.Command{
	Name:  "generate",
	Usage: "gather region's data and into a generated directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "region-id",
			Aliases: []string{"id"},
			Usage:   "region whose data directory will be generated",
		},
		&cli.IntFlag{
			Name:    "quality",
			Aliases: []string{"q"},
			Value:   models.Compressed,
			Usage:   "quality of photos in the datafile (1 - compressed, 2 - original)",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "print extensive logs",
		},
	},
	Action: func(c *cli.Context) error {
		regionID := c.String("region-id")

		quality := models.Quality(c.Int("quality"))
		verbose := c.Bool("verbose")

		if regionID == "" {
			return fmt.Errorf("region id is empty")
		}

		err := generate.Generate(regionID, quality, verbose)
		return err
	},
}

var compressCommand = cli.Command{
	Name:  "compress",
	Usage: "make a zip archive from a generated region directory",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "region-id",
			Aliases: []string{"id"},
			Value:   "",
			Usage:   "region whose generated directory will be compressed",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   false,
			Usage:   "print extensive logs",
		},
	},
	Action: func(c *cli.Context) error {
		regionID := c.String("region-id")
		verbose := c.Bool("verbose")

		if regionID == "" {
			return fmt.Errorf("region id is empty")
		}

		compress.Compress(regionID, verbose)

		return nil
	},
}

var uploadCommand = cli.Command{
	Name:  "upload",
	Usage: "upload a zip archive to the server",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "region-id",
			Aliases: []string{"id"},
			Usage:   "region whose zip archive will be uploaded",
		},

		&cli.IntFlag{
			Name:    "position",
			Aliases: []string{"pos"},
			Usage:   "position at which the datafile will be shown in the app",
		},
		&cli.BoolFlag{
			Name:  "only-meta",
			Usage: "upload only region's metadata, not the zip archive",
		},
		&cli.BoolFlag{
			Name:  "prod",
			Usage: "(dangerous!) upload to production collection (default is test collection)",
		},
	},
	Action: func(c *cli.Context) error {
		regionID := c.String("region-id")

		position := c.Int("position")
		onlyMeta := c.Bool("only-meta")
		prod := c.Bool("prod")

		if regionID == "" {
			return fmt.Errorf("region id is empty")
		}

		if position == 0 {
			return fmt.Errorf("position is 0")
		}

		err := upload.InitFirebase()
		if err != nil {
			return fmt.Errorf("init firebase: %v", err)
		}

		err = upload.Upload(regionID, position, onlyMeta, prod)
		if err != nil {
			return fmt.Errorf("upload %s: %v", regionID, err)
		}

		return nil
	},
}

var optimizeCommand = cli.Command{
	Name:  "optimize",
	Usage: "generate optimized images for a particular place",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "no-icons",
			Value: false,
			Usage: "don't optimize icons",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   false,
			Usage:   "print extensive logs",
		},
	},
	Action: func(c *cli.Context) error {
		noIcons := c.Bool("no-icons")
		verbose := c.Bool("verbose")

		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current working directory: %v", err)
		}

		placeID := filepath.Base(currentDir)

		err = optimize.Optimize(placeID, noIcons, verbose)
		if err != nil {
			return fmt.Errorf("%s: %v", placeID, err)
		}

		return nil
	},
}

func main() {
	app := &cli.App{
		Name:  "touristdb",
		Usage: "manage the tourist database",
		Commands: []*cli.Command{
			&generateCommand,
			&compressCommand,
			&uploadCommand,
			&optimizeCommand,
		},
		CommandNotFound: func(c *cli.Context, command string) {
			log.Printf("invalid command '%s'. See 'touristdb --help'\n", command)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
