package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type branch struct {
	name string
	date string
}

var dryRun bool
var quiet bool
var maxAge int

func main() {

	app := &cli.App{
		Name:  "git-cleanup",
		Usage: "Delete all stale branches",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "dry-run",
				Aliases:     []string{"d"},
				Usage:       "Dry run to find which branches would be deleted",
				Value:       false,
				Destination: &dryRun,
			},
			&cli.BoolFlag{
				Name:        "quiet",
				Aliases:     []string{"q"},
				Usage:       "Hide output",
				Value:       false,
				Destination: &quiet,
			},
			&cli.IntFlag{
				Name:        "age",
				Aliases:     []string{"a"},
				Usage:       "Minimum number of days since last activity to classify a branch as stale",
				Value:       30,
				Destination: &maxAge,
			},
		},
		Action: func(c *cli.Context) (e error) {
			allBranches := make(chan branch)
			staleBranches := make(chan string)

			go GetBranchDates(allBranches)
			go filterForActiveBranches(allBranches, staleBranches)
			go DeleteStaleBranches(staleBranches)

			return
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func filterForActiveBranches(input chan branch, output chan string) {
	defer close(output)
	maxAge *= 24

	for branch := range input {
		lastUpdated, _ := time.Parse("2006-01-02", branch.date)

		if time.Now().Sub(lastUpdated).Hours() > float64(maxAge) {
			if dryRun {
				fmt.Printf("Branch %s is stale", branch.name)
				continue
			}

			output <- branch.name
		}
	}
}

func DeleteStaleBranches(input chan string) {
	for branch := range input {
		err := DeleteBranch(branch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured when attempted to delete branch: %s\n%s\n", branch, err.Error())
			continue
		}

		if !quiet {
			fmt.Printf("Deleted branch: %s", branch)
		}
	}
}
