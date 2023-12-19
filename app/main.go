package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/glaslos/bdiff"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "bdiff",
		Usage: "a tool to create binary diffs",
		Commands: []*cli.Command{
			{
				Name:  "fingerprint",
				Usage: "create a fingerprint for the input file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "Create fingerprint from `FILE`",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fh, err := os.Open(cmd.String("file"))
					if err != nil {
						return err
					}
					fp, err := bdiff.NewFingerprint(fh, 24)
					if err != nil {
						return err
					}
					fp.Print()
					return nil
				},
			},
			{
				Name:  "diff",
				Usage: "create a diff from two files or a fingerprint",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "source",
						Aliases: []string{"s"},
						Usage:   "Diff from `FILE`",
					},
					&cli.StringFlag{
						Name:    "destination",
						Aliases: []string{"d"},
						Usage:   "Diff to `FILE`",
					},
				},
				Action: func(context.Context, *cli.Command) error {
					fmt.Println("boom! I say!")
					return nil
				},
			},
			{
				Name:  "patch",
				Usage: "patch the target file using a diff",
				Action: func(context.Context, *cli.Command) error {
					fmt.Println("boom! I say!")
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
