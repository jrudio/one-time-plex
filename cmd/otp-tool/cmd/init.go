package cmd

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jrudio/go-plex-client"
)

var PlexConn *plex.Plex

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

func initPlex() {
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

	var plexHost string
	var plexToken string

	// get config values

	// read from config
	toml.DecodeFile(cfgFile, &config)
	// _, err := toml.DecodeFile(cfgFile, &config)

	// if err != nil {
	// 	fmt.Println("config decode: ", err)
	// }

	if config.Plex.Host != "" {
		plexHost = config.Plex.Host
		plexToken = config.Plex.Token
	} else {
		plexHost = os.Getenv("PLEXHOST")
		plexToken = os.Getenv("PLEXTOKEN")
	}

	// var err error
	// PlexConn, err = plex.New(plexHost, plexToken)
	PlexConn, _ = plex.New(plexHost, plexToken)

	// if err != nil {
	// 	fmt.Println("plex init: ", err)
	// 	os.Exit(1)
	// }
}
