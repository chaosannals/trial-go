package util

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	outBuf  *bytes.Buffer
	inBuf   io.WriteCloser
	inChan  chan string
	outChan chan string

	client  *ssh.Client
	session *ssh.Session
	user    string
	isStart bool
}

func NewSshClient() *SshClient {
	return &SshClient{
		inBuf:   nil,
		outBuf:  bytes.NewBuffer(make([]byte, 0)),
		inChan:  make(chan string, 1),
		outChan: make(chan string, 1),
		client:  nil,
		session: nil,
		user:    "",
		isStart: false,
	}
}

func (i *SshClient) InitSession(host string, port int, user string, pass string) error {
	fmt.Println("init session")
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(pass))
	clientConfig := &ssh.ClientConfig{
		User:    user,
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
	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return err
	}
	i.client = client
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	i.session = session
	i.user = user
	return nil
}

func (i *SshClient) InitTerminal(host string, port int, user string, pass string) error {
	fmt.Println("init terminal")
	if i.isStart {
		return fmt.Errorf("终端已打开")
	}
	if err := i.InitSession(host, port, user, pass); err != nil {
		return err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := i.session.RequestPty("xterm", 80, 40, modes); err != nil {
		return err
	}
	stdinBuf, err := i.session.StdinPipe()
	if err != nil {
		return err
	}
	i.inBuf = stdinBuf
	i.session.Stdout = i.outBuf

	err = i.session.Shell()
	if err != nil {
		return err
	}

	// io 转发
	ch := make(chan struct{})
	go func() {
		fmt.Println("io transfer")
		buf := make([]byte, 8192)
		var term byte
		if i.user == "root" {
			term = '#'
		} else {
			term = '$'
		}
		var t int
		for {
			if i.outBuf != nil {
				n, err := i.outBuf.Read(buf)
				// _, err := i.outBuf.Read(buf)
				if err != nil && err != io.EOF {
					fmt.Printf("读取错误：%v", err)
					break
				}
				if n > 0 {
					fmt.Println(string(buf[:n]))
				}
				t = bytes.LastIndexByte(buf, term)
				if t > 0 {
					ch <- struct{}{}
					break
				}
			}
		}
	}()
	<-ch
	i.isStart = true
	// fmt.Println("io wait start")
	go i.session.Wait()
	// fmt.Println("io wait end")
	return nil
}

func (i *SshClient) listenMessages() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		cmd := <-i.inChan
		i.inBuf.Write([]byte(fmt.Sprintf("%v\n", cmd)))
		// fmt.Printf("send cmd %v\n", cmd)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		buf := make([]byte, 8192)
		var t int
		term := '$'
		if i.user == "root" {
			term = '#'
		}
		for {
			time.Sleep(time.Millisecond * 200)
			n, err := i.outBuf.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("读 ：%v", err)
				break
			}
			if n > 0 {
				t = bytes.LastIndexByte(buf, byte(term))
				if t > 0 {
					i.outChan <- string(buf[:t])
					break
				}
			} else {
				i.outChan <- string(buf)
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (i *SshClient) SendCmd(cmd string) (string, error) {
	i.inChan <- cmd
	err := i.listenMessages()
	out := <-i.outChan
	return strings.TrimSpace(strings.Split(out, "["+i.user)[0]), err
}
