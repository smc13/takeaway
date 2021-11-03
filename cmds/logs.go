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

	"github.com/fatih/color"
	"github.com/smcassar/takeaway/docker"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Display container logs",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logs(args[0])
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func logs(containerId string) {
	containers, err := docker.StoppableTakeawayContainers()
	if err != nil {
		fmt.Println(err)
		return
	}

	container := docker.FindContainer(containerId, containers)
	if container == nil {
		color.Red("Container %s not found", containerId)
		return
	}

	docker.LogContainer(container)
}
