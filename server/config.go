package main

import (
	"bytes"
	"errors"
	"path/filepath"

	"io/ioutil"

	"os"

	"github.com/BurntSushi/toml"
)

const configFilename = "otpconfig.toml"

// Config holds values that help sets up one time plex
type Config struct {
	Host      string `toml:"host"`
	PlexHost  string `toml:"plexHost"`
	PlexToken string `toml:"plexToken"`
}

// ReadConfig will read a toml file at the given path and help setup or application
func ReadConfig(path string) (Config, error) {
	if path == "" {
		return Config{}, errors.New("path needed for config")
	}

	path = filepath.Join(path, configFilename)

	var config Config

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return config, err
	}

	return config, nil
}

// WriteDefaultConfig creates a configuration file with defaults
func WriteDefaultConfig(path string) (Config, error) {
	config := Config{
		Host:      ":8080",
		PlexHost:  "http://192.168.2.1:32400",
		PlexToken: "abc123",
	}

	buff := bytes.Buffer{}

	if err := toml.NewEncoder(&buff).Encode(config); err != nil {
		return config, err
	}

	if path == "./" {
		cwd, err := os.Getwd()

		if err != nil {
			return config, err
		}

		path = cwd
	}

	path = filepath.Join(path, configFilename)

	if err := ioutil.WriteFile(path, buff.Bytes(), 0644); err != nil {
		return config, err
	}

	return config, nil
}
