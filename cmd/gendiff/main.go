// gendiff is a command-line tool for generating differences between two files.
// It supports various formats including JSON, YAML, and plain text output.
package main

import (
	"code"
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

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
				Usage:       "Format type (stylish|plain|json)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Нужно указать пути до двух файлов
			if cmd.NArg() != 2 {
				err := cli.ShowAppHelp(cmd)

				if err != nil {
					log.Fatal(err)
				}
				return cli.Exit("Error: requires two arguments - path1 and path2 to files", 1)
			}
			pathBefore := cmd.Args().Get(0)
			pathAfter := cmd.Args().Get(1)
			// Валидация format
			allowedFormats := []string{"stylish", "plain", "json"}
			if !slices.Contains(allowedFormats, format) {
				return fmt.Errorf("invalid format '%s'. Must be one of: %s",
					format, strings.Join(allowedFormats, ", "))
			}
			// Получение текста для вывода
			out, err := code.GenDiff(pathBefore, pathAfter, format)
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
