/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmds

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/smcassar/takeaway/docker"
	"github.com/spf13/cobra"
)

var disableAll bool
var removeVolumes bool

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable services",
	Long:  ``,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		containers, err := docker.TakeawayContainers()
		if err != nil {
			color.Red("Error while fetching services: %s", err)
			os.Exit(1)
		}

		services, err := disableableServices(containers)
		if err != nil {
			color.Red("Error while fetching services: %s", err)
			os.Exit(1)
		}

		if disableAll {
			for id := range services {
				disableServiceById(id, containers)
			}

			os.Exit(0)
		}

		if len(services) == 0 {
			color.Cyan("No services to disable")
		}

		if len(args) != 0 {
			for _, arg := range args {
				disableServiceByName(arg, containers, services)
			}

			os.Exit(0)
		}

		selected := showMenu(services)

		if selected != "" {
			disableServiceByName(selected, containers, services)
		}
	},
}

func init() {
	rootCmd.AddCommand(disableCmd)

	disableCmd.Flags().BoolVarP(&disableAll, "all", "a", false, "Disable all service")
	disableCmd.Flags().BoolVarP(&removeVolumes, "volumes", "V", false, "Prune any attached volumes with the service")
}

func showMenu(services map[string]string) string {
	options := make([]string, 0, len(services))
	for _, name := range services {
		options = append(options, name)
	}

	var selected string
	survey.AskOne(&survey.Select{
		Message: "Select services to disable:",
		Options: options,
	}, &selected, survey.WithPageSize(20), SelectIcons)

	return selected
}

func disableableServices(containers []*docker.Container) (map[string]string, error) {
	services := make(map[string]string)

	for _, container := range containers {
		services[container.Id] = strings.Replace(container.Name, "TO--", "", 1)
	}

	return services, nil
}

func disableServiceById(id string, containers []*docker.Container) {
	volumeName, _ := docker.GetAttachedVolumeName(id)
	container := docker.FindContainer(id, containers)

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("fgCyan")
	s.Suffix = color.CyanString(" Disabling Service...")

	s.Start()
	err := docker.RemoveContainer(container)
	s.Stop()
	if err != nil {
		color.Red("Error while disabling service: %s", err)
		return
	}

	color.Cyan("\n✔ Service disabled")

	if volumeName == "" {
		return
	}

	if removeVolumes {
		fmt.Println()

		s.Suffix = color.CyanString(" Removing volume \"%s\"...", volumeName)
		s.Start()
		err = docker.RemoveVolume(volumeName)
		s.Stop()
		if err != nil {
			color.Red("Error while removing volume: %s", err)
			return
		}

		color.Cyan("\n✔ Volume removed")
		return
	}

	color.Yellow(fmt.Sprintf("\nThe disabled service was using a volume named \"%s\"", volumeName))
	color.Yellow(fmt.Sprintf("If you would like to remove this data, run: %s", color.CyanString(fmt.Sprintf("docker volume rm %s", volumeName))))
}

func disableServiceByName(name string, containers []*docker.Container, disableableContainers map[string]string) {
	matches := make(map[string]string)
	for id, serviceName := range disableableContainers {
		if len(serviceName) >= len(name) && serviceName[:len(name)] == name {
			matches[id] = serviceName
		}
	}

	var selectedId string
	switch len(matches) {
	case 0:
		color.Red("No service found with name \"%s\"", name)
		return
	case 1:
		for id := range matches {
			selectedId = id
			break
		}
	default:
		selected := showMenu(matches)
		for id, serviceName := range disableableContainers {
			if serviceName[0:len(selected)] == selected {
				selectedId = id
				break
			}
		}
	}

	disableServiceById(selectedId, containers)
}
