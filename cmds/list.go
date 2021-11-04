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
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/smcassar/takeaway/docker"
	"github.com/spf13/cobra"
)

var asJson bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services enabled by Takeaway",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		containers, err := docker.TakeawayContainers()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if asJson {
			// print containers as json
			out, err := json.Marshal(containers)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(string(out))
			return
		}

		if len(containers) == 0 {
			color.Yellow("No services enabled by Takeaway")
			os.Exit(0)
		}

		fmt.Println("")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Status", "Ports", "Base Alias", "Full Alias"})
		table.SetHeaderColor(tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.FgCyanColor})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(false)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t    ") // pad with tabs
		table.SetNoWhiteSpace(true)

		for _, container := range containers {
			table.Append([]string{container.Id, container.Name, container.Status, container.Ports, container.BaseAlias, container.FullAlias})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&asJson, "json", "j", false, "Output in JSON format")
}
