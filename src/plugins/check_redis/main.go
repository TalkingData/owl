package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	host     string
	password string
	port     string
)

func main() {
	app := cli.NewApp()
	app.Name = "check_redis"
	app.Version = "0.1"
	app.Usage = "redis metric collector"
	app.Authors = []cli.Author{
		{
			Name:  "yingsong",
			Email: "wyingsong@163.com",
		},
	}
	app.Copyright = "2011 (c) TalkingData.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Value:       "127.0.0.1",
			Usage:       "Connect to host.",
			EnvVar:      "redis_host",
			Destination: &host,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "Password to use when connecting to server.",
			EnvVar:      "redis_password",
			Destination: &password,
		},
		cli.StringFlag{
			Name:        "port",
			Value:       "6379",
			Usage:       "Port number to use for connection",
			EnvVar:      "redis_port",
			Destination: &port,
		},
	}
	app.Action = FetchData
	app.Run(os.Args)
}
