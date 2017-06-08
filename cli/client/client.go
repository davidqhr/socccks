package main

import (
	"os"

	"github.com/davidqhr/socccks/client"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ss-local"
	app.Usage = "A go port ShadowSock client"
	app.Version = ""

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "status",
			Usage: "show client running status",
			Action: func(c *cli.Context) error {
				println("not implement yet")
				return nil
			},
		},
	}

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

		client.Start("localhost:8111")

		return nil
	}

	app.Run(os.Args)
}
