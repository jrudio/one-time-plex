package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func checkPlexCredentials() {
	if PlexConn.URL == "" {
		fmt.Println("This command requires Plex credentials")
		fmt.Println()

		fmt.Println("Either provide a `config.toml` file in the current directory or use environment vars PLEXHOST and if needed PLEXTOKEN")

		os.Exit(1)
	}
}

func writeDefaultConfig(filePath string) error {
	if filePath == "" {
		return errors.New("config path is not valid")
	}

	defaultConfig := otpConfig{
		Plex: plexConfig{
			Host:  "http://192.168.1.200:5050",
			Token: "abc123",
		},
	}

	defaultConfigBytes, err := defaultConfig.toBytes()

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, defaultConfigBytes, 0644)
}
