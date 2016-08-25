package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "otp"
	app.Usage = "Automate tasks related to your Plex Media Server"
	app.Action = func(c *cli.Context) error {
		fmt.Println("You have ran otp!")
		return nil
	}

	app.Run(os.Args)
}
