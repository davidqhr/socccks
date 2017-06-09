package main

import (
	"os"

	"github.com/davidqhr/socccks/client"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "socccks-local"
	app.Usage = "Separated Encrypted socks5 proxy client execution"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server, s",
			Usage: "remote server address",
		},
		cli.StringFlag{
			Name:  "local, l",
			Usage: "local socks5 address",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "password to connect remote server",
		},
		cli.BoolFlag{
			Name:  "daemon, d",
			Usage: "run client as daemon a process",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NumFlags() <= 0 {
			cli.ShowAppHelp(c)
			return nil
		}

		// daemonMode := c.GlobalBool("daemon")
		server := c.GlobalString("server")
		local := c.GlobalString("local")
		password := c.GlobalString("password")

		if server == "" {
			return cli.NewExitError("need server address", 1)
		}

		if local == "" {
			return cli.NewExitError("need local address", 1)
		}

		if password == "" {
			return cli.NewExitError("need password", 1)
		}

		client.Start(local, server, password)

		return nil
	}

	app.Run(os.Args)
}
