package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func execCommand(command string, params ...string) error {
	cmd := exec.Command(command, params...)
	cmd.Stderr = cmd.Stdout
	log.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		log.Println(line)
	}
	cmd.Wait()
	return nil
}

func init() {
	logf, err := rotatelogs.New(
		"run/logs/%Y%m%d.log",
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logf)
}

func main() {
	root, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
		return
	}
	jreRoot := filepath.Join(root, "jre")
	jreBin := filepath.Join(jreRoot, "bin")
	envPath := os.Getenv("PATH")
	os.Setenv("PATH", jreBin+";"+envPath)
	if err = execCommand("java", "-version"); err != nil {
		log.Errorln(err)
	}
}
