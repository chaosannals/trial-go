package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chaosannals/trial-go-sshclient/util"
	"gopkg.in/ini.v1"
)

func main() {
	fmt.Println("start")
	cfg, err := ini.Load("ssh.ini")
	if err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}
	host := cfg.Section("").Key("ssh_host").String()
	port, err := cfg.Section("").Key("ssh_port").Int()
	if err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}
	user := cfg.Section("").Key("ssh_user").String()
	pass := cfg.Section("").Key("ssh_password").String()
	client := util.NewSshClient()
	err = client.InitTerminal(host, port, user, pass)
	if err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}
	inReader := bufio.NewReader(os.Stdin)
	for {
		cmd, err := inReader.ReadString('\n')
		if err != nil {
			fmt.Printf("err: %v", err)
		} else {
			msg, err := client.SendCmd(cmd)
			if err != nil {
				fmt.Printf("err: %v", err)
			}
			fmt.Printf("n: %s", msg)
		}
	}
}
