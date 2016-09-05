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

	"github.com/jrudio/go-plex-client"
	"github.com/spf13/cobra"
)

// libraryCmd represents the library command
var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Get a list of your libraries on your Plex server",
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

const libraryTypeOptions = "[movie, show, music, photo, homevideo]"

var newLibraryCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new library. Optionally, create a library just for your Plex friend",
	Long:  "Example: library new /home/jrudio-guest/media/tv --type=show",
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

		// path arg
		argsCount := len(args)
		if argsCount == 0 {
			fmt.Println("one arg needed!")
			fmt.Println("library new [flags] <path-to-library>")
			return
		}

		if argsCount > 1 {
			fmt.Println("only one arg needed!")
			fmt.Println("library new <path-to-library> [flags]")
			return
		}

		if name == "" {
			fmt.Println("'--name' is required")
			cmd.Usage()
			return
		}

		if libraryType == "" {
			fmt.Println("'--type' is required")
			cmd.Usage()
			return
		}

		libraryParams, err := plex.LibraryParamsFromMediaType(libraryType)

		if err != nil {
			fmt.Println(err)
			fmt.Println("flag '--type' must be one of: " + libraryTypeOptions)
			return
		}

		libraryParams.Location = args[0]
		libraryParams.Name = name

		// user can override params with flags
		if libraryParams, err = overrideParams(libraryParams, cmd); err != nil {
			fmt.Println(err)
			return
		}

		// fmt.Printf("%+v\n", libraryParams)

		if shared {
			fmt.Println("creating a shared library")
		} else {
			fmt.Println("creating a new library")
			if err = PlexConn.CreateLibrary(libraryParams); err != nil {
				fmt.Println("create library:", err)
			}
		}

	},
}

var removeLibraryCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a library from your Plex server",
	Run: func(cmd *cobra.Command, args []string) {
		checkPlexCredentials()

		if len(args) != 1 {
			cmd.Usage()
			return
		}

		key := args[0]

		if err := PlexConn.DeleteLibrary(key); err != nil {
			fmt.Printf("failed to remove library: %v\n", err)
			return
		}

		fmt.Printf("Successfully removed %s from your server\n", key)
	},
}

func overrideParams(params plex.CreateLibraryParams, cmd *cobra.Command) (plex.CreateLibraryParams, error) {
	agent, err := cmd.LocalFlags().GetString("agent")
	if err != nil {
		fmt.Println("flag error:", err)
		return plex.CreateLibraryParams{}, err
	}

	scanner, err := cmd.LocalFlags().GetString("scanner")
	if err != nil {
		fmt.Println("flag error:", err)
		return plex.CreateLibraryParams{}, err
	}

	language, err := cmd.LocalFlags().GetString("language")
	if err != nil {
		fmt.Println("flag error:", err)
		return plex.CreateLibraryParams{}, err
	}

	if agent != "" {
		params.Agent = agent
	}

	if scanner != "" {
		params.Scanner = scanner
	}

	if language != "" {
		params.Language = language
	}

	return params, nil
}

func init() {
	RootCmd.AddCommand(libraryCmd)
	libraryCmd.AddCommand(newLibraryCmd)
	libraryCmd.AddCommand(removeLibraryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// libraryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// libraryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newLibraryCmd.Flags().BoolP("shared", "s", false, "Create a shared library")

	newLibraryCmd.Flags().String("name", "", "Name of new library `REQUIRED`")
	newLibraryCmd.Flags().String("type", "", "Library type; One of: "+libraryTypeOptions+" `REQUIRED`")
	newLibraryCmd.Flags().String("agent", "", "Media agent to use to gather metadata for library")
	newLibraryCmd.Flags().String("scanner", "", "Scanner for plex to use")
	newLibraryCmd.Flags().String("language", "en", "Library language `defaults to 'en' (english)")
}
