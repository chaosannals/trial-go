package main

import (
	"os/exec"
	"syscall"
	"os"
	"log"

	/// 此包被加入 work 空间后，会优先拉取 work 指定目录的包
	"github.com/chaosannals/libdemo/simple"
)

/// 启动一个 sh 并进入交互模式，只能在 linux 下执行，调用了 linux 系统方法。要 root
func main() {
	simple.DoSome()
	cmd := exec.Command("/bin/sh")

	// syscall.CLONE_NEWUTS 克隆出新的 UTS 此时 hostname 是隔离的
	// syscall.CLONE_NEWIPC 克隆出新的 IPC
	// syscall.CLONE_NEWPID PID 独立 此时 sh 内部查看自身 PID 为 1
	// syscall.CLONE_NEWNS Mount 名字空间化， mount 和 unmount 挂载的空间独立
	// syscall.CLONE_NEWUSER 用户独立，加上后 sh 用户从 root 变成 nobody
	// syscall.CLONE_NEWNET 网络
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
		// Setpgid: true,
        // Pgid: 0,
		GidMappingsEnableSetgroups: true,

		// 设置了 syscall.CLONE_NEWUSER 这个一设置就有权限问题
		// Credential: &syscall.Credential {
		// 	Uid: uint32(1),
		// 	Gid: uint32(1),
		// },
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}