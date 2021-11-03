package docker

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

const NETWORK_NAME = "takeout"

type Network struct {
	ID   string
	Name string
}

func GetNetworkSettings(alias string, imageName string) []string {
	networkSettings := []string{
		"--network=takeout",
		fmt.Sprintf("--network-alias=%s", alias),
		fmt.Sprintf("--label=com.tighten.takeout.Full_Alias=%s", alias),
	}

	if !baseAliasExists(imageName) {
		networkSettings = append(networkSettings, []string{"--network-alias=" + imageName}...)
		networkSettings = append(networkSettings, []string{fmt.Sprintf("--label=com.tighten.takeout.Base_Alias=%s", imageName)}...)
	}

	return networkSettings
}

func baseAliasExists(name string) bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("\"label=com.tighten.takeout.Base_Alias=%s", name), "--format", "{{.ID}}|{{.Names}}")

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(out) > 0
}

func EnsureNetworkCreated() error {
	networks, err := GetMatchingNetworks()
	if err != nil {
		return err
	}

	if len(networks) == 0 {
		color.Green("Creating missing network: %s", NETWORK_NAME)
		cmd := exec.Command("docker", "network", "create", "-d", "bridge", NETWORK_NAME)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func GetMatchingNetworks() ([]Network, error) {
	cmd := exec.Command("docker", "network", "ls", "--filter", fmt.Sprintf("name=%s", NETWORK_NAME), "--format", "{{.ID}}|{{.Name}}")

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(stdout), "\n")
	networks := make([]Network, len(lines)-1)

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		networks = append(networks, Network{ID: parts[0], Name: parts[1]})
	}

	return networks, nil
}
