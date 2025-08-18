package main

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	//	(&cli.Command{}).Run(context.Background(), os.Args)
	var format string
	cmd := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		UsageText: "gendiff [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				Value:       "stylish",
				Destination: &format,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
