package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davidqhr/socccks/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "socccks-server"
	app.Usage = "Separated Encrypted socks5 proxy server execution"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bindaddress, b",
			Usage: "service address interface",
		},
		cli.StringFlag{
			Name:  "port, p",
			Usage: "service port",
		},
		cli.BoolFlag{
			Name:  "daemon, d",
			Usage: "run server as daemon a process",
		},
	}

	app.Action = func(c *cli.Context) error {
		daemonMode := c.GlobalBool("daemon")
		ip := c.GlobalString("bindaddress")
		port := c.GlobalString("port")

		if ip == "" {
			ip = "0.0.0.0"
		}

		if port == "" {
			port = "8112"
		}

		log.Printf("%s started, version: %s\n", app.Name, app.Version)
		server.Start(fmt.Sprintf("%s:%s", ip, port), daemonMode)
		return nil
	}

	app.Run(os.Args)
}
