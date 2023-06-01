package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/chaosannals/mydocker/container"
)

const usage = `my docker usage.`

var runCommand = cli.Command {
	Name: "run",
	Usage: `my docker run.`,
	Flags: []cli.Flag {
		cli.BoolFlag {
			Name: "ti",
			Usage: "enable tty",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		container.Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command {
	Name: "init",
	Usage: "init container",
	Action: func (context * cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)

		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
	return err
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command {
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}