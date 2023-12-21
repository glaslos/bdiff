package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/glaslos/bdiff"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "bdiff",
		Usage: "a tool to create binary diffs and patch files",
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
				Action: func(ctx context.Context, cmd *cli.Command) error {
					src, err := os.Open(cmd.String("source"))
					if err != nil {
						return err
					}
					dst, err := os.Open(cmd.String("destination"))
					if err != nil {
						return err
					}
					fp, err := bdiff.NewFingerprint(src, 128)
					if err != nil {
						return err
					}

					stat, err := os.Stat(cmd.String("source"))
					if err != nil {
						return err
					}
					diff, err := bdiff.Diff(dst, int(stat.Size()), fp)
					if err != nil {
						return err
					}

					diff.Print()

					return nil
				},
			},
			{
				Name:  "patch",
				Usage: "patch the target file using a diff",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "source",
						Aliases: []string{"s"},
						Usage:   "Target `FILE`",
					},
					&cli.StringFlag{
						Name:    "destination",
						Aliases: []string{"d"},
						Usage:   "Patch for target",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					src, err := os.Open(cmd.String("source"))
					if err != nil {
						return err
					}
					stat, err := os.Stat(cmd.String("source"))
					if err != nil {
						return err
					}

					fp, err := bdiff.NewFingerprint(src, 1024)
					if err != nil {
						return err
					}

					if _, err := src.Seek(0, io.SeekStart); err != nil {
						return err
					}

					delta, err := bdiff.Diff(src, int(stat.Size()), fp)
					if err != nil {
						return nil
					}

					if _, err := src.Seek(0, io.SeekStart); err != nil {
						return err
					}

					dst, err := os.OpenFile(cmd.String("destination"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
					if err != nil {
						return err
					}
					return bdiff.Patch(delta, src, dst)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
