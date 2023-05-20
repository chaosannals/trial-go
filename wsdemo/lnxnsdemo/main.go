package main

import (
	"os/exec"
	"syscall"
	"os"
	"log"
)

/// 启动一个 sh 并进入交互模式，只能在 linux 下执行，调用了 linux 系统方法。要 root
func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}