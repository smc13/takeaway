package docker

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

func ForwardOutputToStdOut(cmd *exec.Cmd) error {
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		outScanner := bufio.NewScanner(outPipe)

		for outScanner.Scan() {
			fmt.Printf("[%s] %s\n", color.GreenString("docker"), outScanner.Text())
		}
	}()

	go func() {
		errScanner := bufio.NewScanner(errPipe)

		for errScanner.Scan() {
			fmt.Printf("[%s] %s\n", color.RedString("docker"), errScanner.Text())
		}
	}()

	return nil
}