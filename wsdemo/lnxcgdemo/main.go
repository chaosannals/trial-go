package main

import (
	"os/exec"
	"path"
	"os"
	"fmt"
	"io/ioutil"
	"syscall"
	"strconv"
)

const cgroupMemoryHierarchMount = "/sys/fs/cgroup/memory"

// 测试 执行首次无效，内存占用 200m，之后执行就起效 占用 100m
// 可能是 hierarch 挂载和 测试程序的启动 顺序上有问题。第2次用了第1次的配置，以此类推。
func main() {
	fmt.Printf("current args[0]: %s", os.Args[0])
	fmt.Println()
	if os.Args[0] == "/proc/self/exe" {
		fmt.Printf("current pid %d", syscall.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr {

		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS ,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error", err)
		os.Exit(2)
	} else {
		fmt.Printf("%v", cmd.Process.Pid)
		os.Mkdir(path.Join(cgroupMemoryHierarchMount, "testmemorylimit"), 0755)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
	}
	cmd.Process.Wait()
}