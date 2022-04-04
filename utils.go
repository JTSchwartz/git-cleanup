package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

func GetBranchDates() (string, error) {
	return execGit([]string{"--no-pager", "branch", "-l", "--format=\"%(committerdate:short) | %(refname:short)\""})
}

func DeleteBranch(branch string) (string, error) {
	return execGit([]string{"branch", "-D", branch})
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

	return strings.TrimRight(stdout.String(), "\000\n"), nil
}
