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

// libraryCmd represents the library command
var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Get a list of your libraries on your Plex server",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		checkPlexCredentials()

		libs, err := PlexConn.GetLibraries()

		if err != nil {
			fmt.Println("get libraries:", err)
			return
		}

		libraryCount := len(libs.Children)

		if libraryCount > 1 {
			fmt.Printf("You have %d libraries\n", libraryCount)
		} else if libraryCount == 1 {
			fmt.Println("You have 1 library")
		} else {
			fmt.Println("No libraries available")
		}

		fmt.Println()

		for _, lib := range libs.Children {
			fmt.Printf("Name: %s\nid: %s\nType: %s\n", lib.Title, lib.Key, lib.Type)
			fmt.Println()
		}
	},
}

var newLibraryCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new shared library for your Plex friend",
	Long:  "A configuration file prevents repetitive credential input",
	Run: func(cmd *cobra.Command, args []string) {
		shared, err := cmd.LocalFlags().GetBool("shared")
		if err != nil {
			fmt.Println("flag error:", err)
			return
		}

		name, err := cmd.LocalFlags().GetString("name")
		if err != nil {
			fmt.Println("flag error:", err)
			return
		}

		libraryType, err := cmd.LocalFlags().GetString("type")
		if err != nil {
			fmt.Println("flag error:", err)
			return
		}

		agent, err := cmd.LocalFlags().GetString("agent")
		if err != nil {
			fmt.Println("flag error:", err)
			return
		}

		scanner, err := cmd.LocalFlags().GetString("scanner")
		if err != nil {
			fmt.Println("flag error:", err)
			return
		}

		location, err := cmd.LocalFlags().GetString("location")
		if err != nil {
			fmt.Println("flag error:", err)
			return
		}

		if shared {
			fmt.Println("creating a shared library")
		} else {
			fmt.Println("creating a new library")
			if err = PlexConn.CreateLibrary(name, location, libraryType, agent, scanner); err != nil {
				fmt.Println("create library:", err)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(libraryCmd)
	libraryCmd.AddCommand(newLibraryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// libraryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// libraryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newLibraryCmd.Flags().BoolP("shared", "s", false, "Create a shared library")

	newLibraryCmd.Flags().String("name", "", "Name of new library `REQUIRED`")
	newLibraryCmd.Flags().String("type", "", "Library type `REQUIRED`")
	newLibraryCmd.Flags().String("agent", "", "Media agent to use to gather metadata for library `REQUIRED`")
	newLibraryCmd.Flags().String("scanner", "", "Scanner for plex to use `REQUIRED`")
	newLibraryCmd.Flags().String("location", "", "Location of the new library `REQUIRED`")
}
