package docker

import (
	"fmt"
	"os/exec"
)

func ImageExists(organization string, image string, tag string) bool {
	cmd := exec.Command("docker", "image", "inspect", fmt.Sprintf("%s/%s:%s", organization, image, tag))

	_, err := cmd.Output()
	return err == nil && cmd.ProcessState.ExitCode() == 0
}

func PullImage(organization string, image string, tag string) error {
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s/%s:%s", organization, image, tag))

	err := ForwardOutputToStdOut(cmd)
	if err != nil {
		return err
	}

	cmd.Run()

	cmd.Wait()

	return nil
}