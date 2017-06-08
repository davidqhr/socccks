package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
			Name:  "config, c",
			Usage: "config path",
		},
	}

	app.Action = func(c *cli.Context) error {
		configPath := c.GlobalString("config")

		if configPath == "" {
			log.Fatalln("no config file")
		}

		config := loadConfig(configPath)

		if config.Address == "" {
			config.Address = "0.0.0.0"
		}

		log.Printf("%s started, version: %s\n", app.Name, app.Version)

		for password, port := range config.Users {
			go server.Start(
				fmt.Sprintf("%s:%d", config.Address, port),
				password)
		}

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGQUIT, syscall.SIGINT)

		<-quit
		log.Println("Quiting...")

		return nil
	}

	app.Run(os.Args)
}
