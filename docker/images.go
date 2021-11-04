package docker

import (
	"fmt"
	"os/exec"
)

func BuildImageName(organization string, imageName string, tag string) string {
	return fmt.Sprintf("%s/%s:%s", organization, imageName, tag)
}

func ImageExists(imgName string) bool {
	cmd := exec.Command("docker", "image", "inspect", imgName)

	_, err := cmd.Output()
	return err == nil && cmd.ProcessState.ExitCode() == 0
}

func PullImage(imgName string) error {
	cmd := exec.Command("docker", "pull", imgName)

	err := ForwardOutputToStdOut(cmd)
	if err != nil {
		return err
	}

	cmd.Run()

	cmd.Wait()

	return nil
}