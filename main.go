package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "git-repo",
		Usage: "Open the current repo in browser",
		Action: func(c *cli.Context) (err error) {
			maxAge := 1
			branches, err := GetBranchDates()

			for _, staleBranch := range filterForActiveBranches(branches, maxAge) {
				fmt.Println(staleBranch)
				err = DeleteBranch(staleBranch)
				if err != nil {
					return err
				}
			}
			return
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
		lastUpdated, _ := time.Parse("2006-10-25", date)
		if time.Now().Sub(lastUpdated).Hours() > float64(maxAge) {
			staleBranches = append(staleBranches, branch)
		}
	}

	return staleBranches
}
