package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	host     string
	user     string
	password string
	port     int
)

func main() {
	app := cli.NewApp()
	app.Name = "check_mysql"
	app.Version = "0.1"
	app.Usage = "mysql performance metric collector"
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
			EnvVar:      "mysql_host",
			Destination: &host,
		},
		cli.StringFlag{
			Name:        "user",
			Value:       "root",
			Usage:       "User for login if not current user",
			EnvVar:      "mysql_user",
			Destination: &user,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "Password to use when connecting to server.",
			EnvVar:      "mysql_password",
			Destination: &password,
		},
		cli.IntFlag{
			Name:        "port",
			Value:       3306,
			Usage:       "Port number to use for connection",
			EnvVar:      "mysql_port",
			Destination: &port,
		},
	}
	app.Action = FetchData
	app.Run(os.Args)
}
