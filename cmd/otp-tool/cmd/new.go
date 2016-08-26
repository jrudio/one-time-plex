// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a config file, or create a shared library",
	// Long: ``,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// TODO: Work your own magic here
	// 	fmt.Println("new called")
	// },
}

var defaultConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate a default configuration file",
	Long:  "A configuration file prevents repetitive credential input",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.LocalFlags().GetString("path")

		if err != nil {
			fmt.Println("flag error:", err)
			return
		}

		if err := writeDefaultConfig(path); err != nil {
			fmt.Println("write config:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.AddCommand(defaultConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	defaultConfigCmd.Flags().String("path", "./config.toml", "Path to write the default configuration file")

}
