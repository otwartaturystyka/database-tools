package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFlags(0)
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
	},
	Action: func(c *cli.Context) error {
		if c.String("region-id") == "" {
			return errors.New("region-id is empty")
		}

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
			&compressCommand,
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
