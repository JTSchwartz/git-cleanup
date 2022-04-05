package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	var dryRun bool
	var quiet bool
	var maxAge int

	app := &cli.App{
		Name:  "git-repo",
		Usage: "Open the current repo in browser",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "dry-run",
				Aliases:     []string{"d"},
				Usage:       "Dry run to find which branches would be deleted.",
				Value:       false,
				Destination: &dryRun,
			},
			&cli.BoolFlag{
				Name:        "quiet",
				Aliases:     []string{"q"},
				Usage:       "Hide output.",
				Value:       false,
				Destination: &quiet,
			},
			&cli.IntFlag{
				Name:        "age",
				Aliases:     []string{"a"},
				Usage:       "Max age: Any branches n+ days old will be deleted. Default: 30",
				Value:       30,
				Destination: &maxAge,
			},
		},
		Action: func(c *cli.Context) error {
			if branches, err := GetBranchDates(); err == nil {
				for _, staleBranch := range filterForActiveBranches(branches, maxAge) {
					if !quiet {
						fmt.Printf("Deleted branch: %s", staleBranch)
					}

					if dryRun {
						fmt.Printf("Branch %s is stale", staleBranch)
						continue
					}

					err = DeleteBranch(staleBranch)
					if err != nil {
						fmt.Fprintf(os.Stderr, "An error occured while deleting branch: %s\n%s", staleBranch, err.Error())
					}
				}
			} else {
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func filterForActiveBranches(branches map[string]string, maxAge int) []string {
	maxAge *= 24
	staleBranches := make([]string, 0)

	for branch, date := range branches {
		lastUpdated, _ := time.Parse("2006-01-02", date)
		if time.Now().Sub(lastUpdated).Hours() > float64(maxAge) {
			staleBranches = append(staleBranches, branch)
		}
	}

	return staleBranches
}
