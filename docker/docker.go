package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Container struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Ports     string `json:"ports"`
	BaseAlias string `json:"base_alias"`
	FullAlias string `json:"full_alias"`
}

func (c *Container) String() string {
	return c.Name
}

func TakeawayContainers() ([]*Container, error) {
	cmd := exec.Command("docker", "ps", "-a", "--filter", "name=TO-", "--format", "{{.ID}}|{{.Names}}|{{.Status}}|{{.Ports}}|{{.Label \"com.tighten.takeout.Base_Alias\"}}|{{.Label \"com.tighten.takeout.Full_Alias\"}}")

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(string(stdout), err.Error())
		return nil, err
	}

	// get lines from stdout
	lines := strings.Split(string(stdout), "\n")

	// map lines into containers
	containers := make([]*Container, len(lines)-1)
	for i, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		containers[i] = &Container{
			Id:        parts[0],
			Name:      parts[1],
			Status:    parts[2],
			Ports:     parts[3],
			BaseAlias: parts[4],
			FullAlias: parts[5],
		}
	}

	return containers, nil
}

func StartableTakawayContainers() ([]*Container, error) {
	containers, err := TakeawayContainers()
	if err != nil {
		return nil, err
	}

	// filter out containers that are already running
	startableContainers := []*Container{}
	for _, container := range containers {
		if strings.Contains(container.Status, "Up") {
			continue
		}
		startableContainers = append(startableContainers, container)
	}

	return startableContainers, nil
}

func StoppableTakeawayContainers() ([]*Container, error) {
	containers, err := TakeawayContainers()
	if err != nil {
		return nil, err
	}

	// filter out containers that arent already running
	stoppableContainers := []*Container{}
	for _, container := range containers {
		if !strings.Contains(container.Status, "Up") {
			continue
		}
		stoppableContainers = append(stoppableContainers, container)
	}

	return stoppableContainers, nil
}

func FindContainer(id string, containers []*Container) *Container {
	if strings.Contains(id, " -") {
		id = strings.Split(id, " -")[0]
	}

	for _, container := range containers {
		if container.Id == id || container.Name == id {
			return container
		}
	}

	return nil
}

func StartContainer(container *Container) error {
	cmd := exec.Command("docker", "start", container.Id)
	_, err := cmd.Output()
	return err
}

func StopContainer(container *Container) error {
	cmd := exec.Command("docker", "stop", container.Id)
	_, err := cmd.Output()
	return err
}

// forward container logs to stdout
func LogContainer(container *Container) error {
	cmd := exec.Command("docker", "logs", "-f", container.Id)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// create a new container
func CreateContainer(args []string) error {
	cmd := exec.Command("docker", append([]string{"run", "-d"}, args...)...)

	err := ForwardOutputToStdOut(cmd)
	if err != nil {
		return err
	}

	cmd.Start()

	return cmd.Wait()
}

// remove a container
func RemoveContainer(container *Container) error {
	stoppable, err := StoppableTakeawayContainers()
	if err != nil {
		return err
	}

	// if in stoppable list, stop it first
	for _, c := range stoppable {
		if c.Id == container.Id {
			err := StopContainer(c)
			if err != nil {
				return err
			}
		}
	}

	cmd := exec.Command("docker", "rm", container.Id)

	_, err = cmd.Output()

	return err
}
