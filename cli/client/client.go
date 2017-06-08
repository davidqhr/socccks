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
			Usage: "remote shadowsocks server ip",
		},
		cli.StringFlag{
			Name:  "port, p",
			Usage: "remote shadowsocks server port",
		},
		cli.StringFlag{
			Name:  "password, P",
			Usage: "password to connect shadowsocks server",
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

		daemonMode := c.GlobalBool("daemon")
		serverIP := c.GlobalString("server")
		serverPort := c.GlobalString("port")
		password := c.GlobalString("password")

		if serverIP == "" {
			return cli.NewExitError("need server ip", 1)
		}

		if serverPort == "" {
			return cli.NewExitError("need server port", 1)
		}

		if password == "" {
			return cli.NewExitError("need password", 1)
		}

		println(daemonMode, serverIP, serverPort, password)

		client.Start("0.0.0.0:8111")

		return nil
	}

	app.Run(os.Args)
}
