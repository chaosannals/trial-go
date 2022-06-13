package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chaosannals/trial-go-stress/stress"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "scheme",
			Aliases: []string{"s"},
			Value:   "http",
		},
		&cli.StringFlag{
			Name:    "host",
			Aliases: []string{"a"},
			Value:   "127.0.0.1",
			Usage:   "host",
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "80",
			Usage:   "port",
		},
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"d"},
			Value:   "",
			Usage:   "url path",
		},
		&cli.StringFlag{
			Name:    "method",
			Aliases: []string{"m"},
			Value:   "GET",
			Usage:   "http method",
		},
		&cli.StringFlag{
			Name:    "times",
			Aliases: []string{"t"},
			Value:   "100",
			Usage:   "request times",
		},
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Value:   "",
			Usage:   "config file path",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Value:   true,
			Usage:   "show info",
		},
	}

	app.Action = func(c *cli.Context) error {
		cfgPath := c.String("conf")

		if cfgPath == "" {
			port, err := strconv.Atoi(c.String("port"))
			if err != nil {
				return err
			}
			wkr := stress.NewHttpStressWorker(
				&stress.StressConfig{
					Scheme: c.String("scheme"),
					Host:   c.String("host"),
					Port:   port,
					Path:   c.String("path"),
				},
			)
			wkr.Request()
		} else {
			cfg, err := stress.LoadConfig(cfgPath)
			fmt.Printf("load config: %s", cfgPath)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			wkr := stress.NewHttpStressWorker(cfg)
			wkr.Request()
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}
	fmt.Println("final.")
}
