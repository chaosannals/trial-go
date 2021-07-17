package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/ini.v1"
)

func connect() (*ssh.Session, error) {
	cfg, err := ini.Load("ssh.ini")
	if err != nil {
		return nil, err
	}
	host := cfg.Section("").Key("ssh_host").String()
	port, err := cfg.Section("").Key("ssh_port").Int()
	if err != nil {
		return nil, err
	}
	user := cfg.Section("").Key("ssh_user").String()
	pass := cfg.Section("").Key("ssh_password").String()
	key := ""
	fmt.Printf("connect: %s", host)
	auth := make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(pass))
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}
		var signer ssh.Signer
		if pass == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(pass))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	cipherList := make([]string, 0)
	var client *ssh.Client
	var config ssh.Config
	if len(cipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256",
				"arcfour128",
				"aes128-cbc",
				"3des-cbc",
				"aes192-cbc",
				"aes256-cbc",
			},
		}
	} else {
		config = ssh.Config{
			Ciphers: cipherList,
		}
	}

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}
	return session, nil
}

func main() {
	fmt.Println("start")

	session, err := connect()
	if err != nil {
		fmt.Printf("ssh error: %v", err)
		os.Exit(1)
	}
	defer session.Close()
	outbt := bytes.Buffer{}
	errbt := bytes.Buffer{}
	session.Stdout = &outbt
	session.Stderr = &errbt
	err = session.Run("ls -l /usr")
	if err != nil {
		fmt.Printf("ssh error: %v", err)
		os.Exit(2)
	}
	fmt.Println(outbt.String())
	fmt.Println(errbt.String())
}
