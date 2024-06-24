package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	// "path/filepath"

	"os/user"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// homeDir, err := filepath.Abs("~") // Windows 下无法通过 ~ 获取到
	userInfo, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("homeDir: %s\n", userInfo.HomeDir)

	privKeyPath := filepath.Join(userInfo.HomeDir, ".ssh", "id_rsa")
	fmt.Printf("privKey: %s\n", privKeyPath)
	privKeyText, err := os.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("privKey content: \n%s\n", privKeyText)

	// privKeyBlock, _ := pem.Decode(privKeyText)
	signer, err := ssh.ParsePrivateKey(privKeyText)
	if err != nil {
		log.Fatalln(err)
	}

	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.PublicKeys(signer))
	clientConfig := &ssh.ClientConfig{
		User:    "root",
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config: ssh.Config{
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
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	host := os.Getenv("SSH_DEMO_HOST")
	port := os.Getenv("SSH_DEMO_PORT")

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("target: %s\n", addr)

	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		log.Fatalln(err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatalln(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatalln(err)
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
