package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func GetBranchDates(output chan branch) {
	defer close(output)
	branchesByDate, e := execGit([]string{"--no-pager", "branch", "-l", "--format=\"%(committerdate:short)|%(refname:short)\""})
	if e != nil {
		fmt.Fprintln(os.Stderr, "Unable to find any git branches.")
		close(output)
		return
	}

	for _, str := range strings.Split(branchesByDate, "\n") {
		parts := strings.SplitN(str, "|", 2)
		output <- branch{
			name: parts[1],
			date: parts[0],
		}
	}
	return
}

func DeleteBranch(branch string) (err error) {
	_, err = execGit([]string{"branch", "-D", branch})
	return
}

func execGit(gitArgs []string) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", errors.New("wait status returned non-zero")
			}
		}
		return "", err
	}

	output := strings.ReplaceAll(stdout.String(), "\"", "")
	return strings.TrimRight(output, "\000\n"), nil
}
