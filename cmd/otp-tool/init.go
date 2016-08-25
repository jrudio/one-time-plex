package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type otpConfig struct {
	Plex plexConfig `toml:"plex"`
}

type plexConfig struct {
	Host  string `toml:"host"`
	Token string `toml:"token"`
}

func (otp otpConfig) toBytes() ([]byte, error) {
	var b bytes.Buffer

	err := toml.NewEncoder(&b).Encode(otp)

	return b.Bytes(), err
}

var config otpConfig

func init() {
	// TODO change behavior of init
	//
	// check PMS credentials via different options
	// do not force a config file
	//
	// if one is present read from it
	// else if check environment var
	// else ignore
	//
	// if a command requires PMS credentials
	// and they are not present
	// then an error will suffice
	//
	configFilePath := flag.String("config", "", "./config.toml")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "config.toml"
	}

	// get config values
	_, err := toml.DecodeFile(*configFilePath, &config)

	if err != nil {
		// likely file not found error, so write the default config to file
		fmt.Println(err)
		fmt.Println("creating default config file...")

		err = writeDefaultConfig(*configFilePath)

		if err != nil {
			// failing here means we should exit program
			fmt.Println(err)
		}

		fmt.Println("please edit the config.toml file")

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
