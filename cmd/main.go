package main

import (
	"errors"
	"log"
	"os"

	"github.com/bartekpacia/database-tools/cmd/compress"
	"github.com/bartekpacia/database-tools/cmd/generate"
	"github.com/bartekpacia/database-tools/cmd/upload"
	"github.com/bartekpacia/database-tools/models"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFlags(0)
}

var generateCommand = cli.Command{
	Name:  "generate",
	Usage: "gather region's data and put them into a generated directory",
	OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
		log.Println("error:", err)
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "region-id",
			Aliases: []string{"id"},
			Value:   "",
			Usage:   "region whose data directory will be generated",
		},
		&cli.StringFlag{
			Name:  "lang",
			Value: "pl",
			Usage: "language of the generated directory",
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
			Value:   false,
			Usage:   "print extensive logs",
		},
	},
	Action: func(c *cli.Context) error {
		regionID := c.String("region-id")
		lang := c.String("lang")
		quality := models.Quality(c.Int("quality"))
		verbose := c.Bool("verbose")

		if regionID == "" {
			return errors.New("region-id is empty")
		}

		err := generate.Generate(regionID, lang, quality, verbose)
		return err
	},
}

var compressCommand = cli.Command{
	Name:  "compress",
	Usage: "make a zip archive from a generated region directory",
	OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
		log.Println("error:", err)
		return nil
	},
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
			return errors.New("region-id is empty")
		}

		compress.Compress(regionID, verbose)

		return nil
	},
}

var uploadCommand = cli.Command{
	Name:  "upload",
	Usage: "upload a zip archive to the server",
	OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
		log.Println("error:", err)
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "region-id",
			Aliases: []string{"id"},
			Value:   "",
			Usage:   "region whose zip archive will be uploaded",
		},
		&cli.StringFlag{
			Name:  "lang",
			Value: "pl",
			Usage: "language of the zip archive that will be uploaded",
		},
		&cli.IntFlag{
			Name:    "position",
			Aliases: []string{"pos"},
			Value:   1,
			Usage:   "position at which the datafile will be shown in the app",
		},
		&cli.BoolFlag{
			Name:  "only-meta",
			Value: false,
			Usage: "upload only region's metadata, not the zip archive",
		},
		&cli.BoolFlag{
			Name:  "prod",
			Value: false,
			Usage: "(dangerous!) upload to production collection (default is test collection)",
		},
	},
	Action: func(c *cli.Context) error {
		regionID := c.String("region-id")
		lang := c.String("lang")
		position := c.Int("position")
		onlyMeta := c.Bool("onlyMeta")
		prod := c.Bool("prod")

		if regionID == "" {
			return errors.New("region-id is empty")
		}

		upload.Upload(regionID, lang, position, onlyMeta, prod)
		return nil
	},
}

func main() {
	app := &cli.App{
		Name:  "touristdb",
		Usage: "manage the tourist database",
		OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
			log.Println("error:", err)
			return nil
		},
		Commands: []*cli.Command{
			&generateCommand,
			&compressCommand,
			&uploadCommand,
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
