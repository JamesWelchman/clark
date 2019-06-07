/*
pacman counts the number of updates to install

  numUpdates, err := pacman.NumUpdates()

numUpdates is an int
*/
package pacman

import (
	"fmt"
	"os/exec"
	"strings"
)

var cmd = [...]string{
	"/usr/bin/pacman",
	"-Qu",
}

func runPacman() (string, error) {
	writer := &strings.Builder{}

	cmd := exec.Cmd{
		Path:   cmd[0],
		Args:   cmd[:],
		Stdout: writer,
	}
	err := cmd.Run()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if !ok {
			return "", err
		}

		if exitError.ExitCode() == 1 {
			return "", nil
		}

		return "", exitError
	}

	return writer.String(), nil
}

func countUpdates(output string) int {
	total := 0

	for _, r := range output {
		if r == '\n' {
			total++
		}
	}

	return total
}

func NumUpdates() (int, error) {
	output, err := runPacman()
	if err != nil {
		return 0, fmt.Errorf("couldn't run pacman :: %v", err)
	}

	return countUpdates(output), nil
}
