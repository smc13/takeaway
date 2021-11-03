package docker

import (
	"encoding/json"
	"os/exec"
)

type Volumes struct {
	Name string
}

func GetAttachedVolumeName(id string) (string, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{json .Mounts}}", id)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var volumes []Volumes
	err = json.Unmarshal(out, &volumes)
	if err != nil {
		return "", err
	}

	if len(volumes) != 0 {
		return volumes[0].Name, nil
	}

	return "", nil
}

func RemoveVolume(name string) error {
	cmd := exec.Command("docker", "volume", "rm", name)

	_, err := cmd.Output()

	return err
}
