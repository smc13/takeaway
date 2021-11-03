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
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/smcassar/takeaway/services"
	"github.com/spf13/cobra"
)

var useDefaults bool
var parallel bool

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable services",
	Long:  ``,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, arg := range args {
				enableService(arg)
			}

			os.Exit(0)
		}

		var selected []string
		survey.AskOne(&survey.MultiSelect{
			Message: "Select services to enable:",
			Options: services.GetServiceNames(),
		}, &selected, survey.WithPageSize(10), SelectIcons)

		for _, name := range selected {
			enableService(name)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(enableCmd)

	enableCmd.Flags().BoolVarP(&useDefaults, "defaults", "d", false, "Use default configurations")
	enableCmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "Enable services in parallel")
}

func enableService(name string) {
	service := services.GetService(name)
	if service == nil {
		color.Red("Service %s not found", name)
		return
	}

	services.EnableService(service, useDefaults)
}
