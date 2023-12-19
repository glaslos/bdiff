package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
				Action: func(context.Context, *cli.Command) error {
					fmt.Println("boom! I say!")
					return nil
				},
			},
			{
				Name:  "diff",
				Usage: "create a diff from two files or a fingerprint",
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
