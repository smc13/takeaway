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

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/smcassar/takeaway/docker"
	"github.com/spf13/cobra"
)

var stopAll bool

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all containers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		containers, err := docker.StoppableTakeawayContainers()
		if err != nil {
			fmt.Println("Error getting containers:", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			found := 0
			for _, arg := range args {
				// find container by name or id
				container := docker.FindContainer(arg, containers)
				if container == nil {
					color.Red("Container not found: %s", arg)
				} else {
					found++
					stopContainer(container)
				}
			}

			color.Cyan("✔ %d Containers stopped", found)
			os.Exit(0)
		}

		if stopAll {
			for _, container := range containers {
				stopContainer(container)
			}

			color.Cyan("✔ %d Containers stopped", len(containers))
			os.Exit(0)
		}

		if len(containers) == 0 {
			color.Yellow("No containers to stop")
			os.Exit(1)
		}

		var choices []string
		for _, container := range containers {
			choices = append(choices, fmt.Sprintf("%s - %s", container.Id, container.Name))
		}

		var selected []string
		survey.AskOne(&survey.MultiSelect{
			Message: "Select containers to stop",
			Options: choices,
		}, &selected, SelectIcons)

		for _, name := range selected {
			for _, container := range containers {
				if fmt.Sprintf("%s - %s", container.Id, container.Name) == name {
					stopContainer(container)
				}
			}
		}

		color.Cyan("✔ %d Containers stopped", len(selected))
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVarP(&stopAll, "all", "a", false, "stop all containers")
}

func stopContainer(container *docker.Container) {
	if err := docker.StopContainer(container); err != nil {
		color.Red("Error stopping container %s: %s", container.Name, err)
	} else {
		color.Green("✔ Stopped %s", container.Name)
	}
}
