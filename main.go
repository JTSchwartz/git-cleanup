package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "git-repo",
		Usage: "Open the current repo in browser",
		Action: func(c *cli.Context) (err error) {
			branches, err := GetBranchDates()
			fmt.Println(branches)
			return
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
