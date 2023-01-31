package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"godeps/internal/crawler"
	"godeps/internal/presentation"
	"golang.org/x/mod/modfile"
	"io"
	"log"
	"os"
)

const goMod = "go.mod"

type Config struct {
	Crawler crawler.Config
}

func ExportDependencies(root string, outputPath string, config Config) error {
	var out io.Writer = os.Stdout
	if outputPath != "" {
		file, err := os.Open(outputPath)
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		defer file.Close()
	}

	if root == "." {
		mod, err := os.ReadFile(goMod)
		if err != nil {
			return fmt.Errorf("could not open %s", goMod)
		}
		f, err := modfile.ParseLax("go.mod", mod, nil)
		root = f.Module.Mod.String()
	}

	scanner := crawler.Crawler{}
	imports := scanner.Crawl(root, config.Crawler)
	err := presentation.AsDotNotation(imports, out)
	if err != nil {
		return fmt.Errorf("could not generate output: %w", err)
	}
	return nil
}

func main() {
	var cfg Config
	app := &cli.App{
		Name:      "godeps",
		Usage:     "visualize the dependencies in a Go package",
		UsageText: "godeps [global options] <root pkg>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "o",
				Value: "",
				Usage: "filename for output",
			},
			&cli.IntFlag{
				Name:    "depth",
				Aliases: []string{"d"},
				Value:   0,
				Usage:   "max recursion depth",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return errors.New("requires root package argument; see `godeps help`")
			}
			err := ExportDependencies(c.Args().Get(0), c.Args().Get(1), cfg)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed:", err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
