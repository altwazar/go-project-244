package main

import (
	"code"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Нужен один аргумент - путь
			if cmd.NArg() != 2 {
				err := cli.ShowAppHelp(cmd)

				if err != nil {
					log.Fatal(err)
				}
				return cli.Exit("Error: requires two argumenst - path1 and path2 to files", 1)
			}

			pathBefore := cmd.Args().Get(0)
			pathAfter := cmd.Args().Get(1)
			out, err := code.CompareConfigs(pathBefore, pathAfter, format)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
