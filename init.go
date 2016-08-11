package main

import (
	"bytes"
	"flag"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

type otpConfig struct {
	Host                string     `toml:"host"`
	MonitorUserInterval int        `toml:"monitorUserInterval"`
	Plex                plexConfig `toml:"plex"`
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
	configFilePath := flag.String("config", "", "./config.toml")

	flag.Parse()

	if *configFilePath == "" {
		*configFilePath = "config.toml"
	}

	// get config values
	_, err := toml.DecodeFile(*configFilePath, &config)

	if err != nil {
		// likely file not found error, so write the default config to file
		log.Warn(err)
		log.Info("creating default config file...")

		err = writeDefaultConfig()

		if err != nil {
			// failing here means we should exit program
			log.WithError(err).Fatal("failed to write config")
		}

		log.Info("please edit the config.toml file")

		os.Exit(1)
	}
}

func writeDefaultConfig() error {
	defaultConfig := otpConfig{
		Host:                ":4040",
		MonitorUserInterval: 1500,
		Plex: plexConfig{
			Host:  "http://192.168.1.200:5050",
			Token: "abc123",
		},
	}

	defaultConfigBytes, err := defaultConfig.toBytes()

	if err != nil {
		return err
	}

	return ioutil.WriteFile("config.toml", defaultConfigBytes, 0644)
}
