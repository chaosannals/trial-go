package container

import (
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	log "github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	absExe, err := filepath.Abs(os.Args[0])
	if (err != nil) {
		log.Errorf(err.Error())
	}
	log.Infof("NewParentProcess %s", absExe)
	// 这个拿不到
	selfExe, err := os.Readlink("/proc/self/exe")
	if (err != nil) {
		log.Errorf(err.Error())
	}
	args := []string { "init", command }
	cmd := exec.Command(selfExe, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func Run(tty bool, command string) error {
	log.Info("Run")
	parent := NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	return nil
}

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	// 显示 mount namespace 独立
	syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")

	// MS_NOEXEC 不允许运行程序
	// MS_NOSUID 不允许 set-user-ID 或 set-group-ID
	// MS_NODEV 
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{ command }
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}