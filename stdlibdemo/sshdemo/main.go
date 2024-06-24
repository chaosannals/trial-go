package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"os/user"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
)

type SshConn struct {
	user    string
	client  *ssh.Client
	session *ssh.Session
	term    byte

	outBuf *bytes.Buffer
	inBuf  io.WriteCloser

	output []byte // 输出残余
}

func ConnSsh(
	host string,
	port string,
	user string,
	privKeyText []byte,
) (*SshConn, error) {
	// privKeyBlock, _ := pem.Decode(privKeyText)
	signer, err := ssh.ParsePrivateKey(privKeyText)
	if err != nil {
		return nil, err
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

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("target: %s\n", addr)

	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
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

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	outBuf := bytes.NewBuffer(make([]byte, 0))
	session.Stdout = outBuf

	var term byte
	if user == "root" {
		term = '#'
	} else {
		term = '$'
	}

	return &SshConn{
		user:    user,
		client:  client,
		session: session,
		inBuf:   stdinPipe,
		outBuf:  outBuf,
		term:    term,
		output:  make([]byte, 0),
	}, nil
}

func (conn *SshConn) Close() {
	if conn.session != nil {
		conn.session.Close()
	}
	if conn.client != nil {
		conn.client.Close()
	}
}

func (conn *SshConn) Exec(cmds ...string) error {
	err := conn.session.Shell()
	if err != nil {
		return err
	}

	// 开始信息
	_, err = conn.ReadOutput(0)
	if err != nil {
		return err
	}
	// fmt.Print(string(startInfo))

	var output []byte
	for i, cmd := range cmds {
		line := []byte(fmt.Sprintf("%v\n", cmd))
		_, err := conn.inBuf.Write(line)
		if err != nil {
			return err
		}

		all, err := conn.ReadOutput(i + 1)
		if err != nil {
			return err
		}
		output = all
	}

	blocks := bytes.Split(output, []byte{conn.term})
	for i, block := range blocks {
		fmt.Print(string(block))
		if len(block) > 0 {
			fmt.Print(string(conn.term))
		}
		if i < len(cmds) {
			fmt.Println(cmds[i])
		}
	}

	return nil
}

// Shell 模式是整个缓冲全出。
// 所以每次读都是整个屏幕包括上次重复的部分。
// sshd 发送过来的就是如此，客户端无法修改。
func (conn *SshConn) ReadOutput(index int) ([]byte, error) {
	output := append(make([]byte, 0), conn.output...)
	for {
		buf := make([]byte, 8192)
		n, err := conn.outBuf.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if n > 0 {
			output = append(output, buf[:n]...)

			blocks := bytes.Split(output, []byte{conn.term})
			if len(blocks) > index {
				t := bytes.LastIndexByte(output, conn.term)
				if t > 0 {
					output = output[:t+1]
					conn.output = append(make([]byte, 0), output[t+1:]...)
					break
				}
			}
		} else {
			time.Sleep(time.Millisecond * 10)
		}
	}
	return output, nil
}

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
	// fmt.Printf("privKey content: \n%s\n", privKeyText)

	host := os.Getenv("SSH_DEMO_HOST")
	port := os.Getenv("SSH_DEMO_PORT")

	conn, err := ConnSsh(host, port, "root", privKeyText)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	fmt.Println("====================================")

	err = conn.Exec(
		"cd ~",
		"pwd",
		"ls -al",
	)
	if err != nil {
		log.Fatalln(err)
	}
}
