# Linux 容器 Demo

## Namespace 名字空间【隔离】

```bash
## 打印当前进程树
pstree -pl

## 查看当前 PID
echo $$

## 打印进程 UTS  444 为 PID
readlink /proc/444/ns/uts

## 查看 hostname
hostname

## 修改 hostname
hostname -b aaaa
```

```bash
## 查看消息队列
ipcs -p

## 查看统计的 IPC 状态
ipcs -u

## 创建 IPC
ipcmk -Q
```

```bash
## 查看系统信息
ls /proc

## 查看进程信息
ps -ef

## 指定类型 -t proc；  当前 proc 挂载到 /proc
mount -t proc proc /proc
```

```bat
@rem 通过 root 启动 wsl
wsl -u root
```

```bash
## 查询当前用户
whoami

## 查看当前用户 id 信息
id

## 查看用户信息文件
cat /etc/passwd
```

```bash
# 查看网络地址
ip address

# 查看网络路由
ip route

# 查看网络(这个比较老的命令，可能新系统没有要装 net-tools)
ifconfig
```

## Cgroups 资源【限制】

附录：类比 Windows 下 JobObject 

- task 进程 或 线程
- subsystem 子系统，具体资源控制器
- cgroup  关联 task 和 subsystem 的表示
- hierarchy 层级树，可以挂多个 subsystem


```bash
# 老版本
apt-get install cgroup-bin

# 新版本
apt install cgroup-tools

# 列举子系统 lssubsys 命令 cgroup-tools 工具集提供
lssubsys -a
```

```bash
# 查看 进程 控制组 信息
cat /proc/cgroups

# 挂载 cgroup
mount -t cgroup -o none,name=cgroup-test cgroup-test ./cgroup-test

# 列举控制组
ls ./cgroup-test

# 到 cgroup 目录下创建的目录也是 cgroup 类型的 子 cgroup
cd ./cgroup-test
mkdir cgroup-1
mkdir cgroup-2

# 删除 子 cgroup
rmdir cgroup-1

# 查看当前 PID
echo $$
# 查看进程 cgroup
cat /proc/[pid]/cgroup

# 把进程挂载到 子 cgroup
sh -c "echo $$ >> ./cgroup-1/tasks"

# 在 cgroup 目录下执行可以查看 cgroup 项
# 查看 memory 项 在 /sys/fs/cgroup/memory
mount | grep memory

# stress 命令用于模拟负载
# -c ; --cpu N 产生 N 个进程，反复计算随机数平方根
# -i ; --io N 产生 N 个进程，反复调用 sync() 将内存随机内容写入硬盘
# -m ; --vm N 产生 N 个进程，不断分配释放内存
#      --vm-bytes B 指定内存大小
#      --vm-stride B 不断给部分内存赋值，让 Copy On Write 发生
#      --vm-hang N 分配到内存后随眠 N 秒
#      --vm-keep  一直占用内存
# 其他自行查阅 -h

# 启动一个 stress 进程 200M 内存，会占用 sh
stress --vm-bytes 200m --vm-keep -m 1

# 打开其他终端用 top 查看内存可以发现此 stress 的 sh 占用 200m 内存。
top

# 关闭 stress 后执行以下

# 在 /sys/fs/cgroup/memory 子系统里面建子 cgroup 并进入
cd /sys/fs/cgroup/memory
mkdir test-limit-memory
cd test-limit-memory

# 写入内存限制
sudo sh -c "echo '100m' > memory.limit_in_bytes"

# 把当前进程挂到 cgroup 上去，此时当前的 sh 被新的 test-limit-memory 限制内存。
sudo sh -c "echo $$ > tasks" 

# 这个好像可以不用，只挂 tasks 即可。
sudo sh -c "echo $$ > cgroup.procs" 

# 再次启动一个 stress 进程 200M 内存，会占用 sh
stress --vm-bytes 200m --vm-keep -m 1

# 打开其他终端用 top 查看内存可以发现此 stress 的 sh 占用 100m 内存。
top
```

## AUFS

```bash
# 创建挂载目录
mkdir mnt
# 创建容器文件
mkdir container-layer &&  echo "container layer." > ./container-layer/container-layer.txt
mkdir image-1 && echo "iamge layer 1" > ./image-1/image-1.txt
mkdir image-2 && echo "iamge layer 2" > ./image-2/image-2.txt
mkdir image-3 && echo "iamge layer 3" > ./image-3/image-3.txt
mkdir image-4 && echo "iamge layer 4" > ./image-4/image-4.txt

# 查看操作系统支持的文件系统
ls /sys/fs/ -al

# 安装 aufs-tools , ubuntu22.04 删了它。
apt update
apt install aufs-tools

# vbox 和 wsl2 都可以通过下面安装
apt install linux-image-extra-virtual
# vbox 可以使用，WSL2 不能上 aufs，下面命令会提示找不到 windows-wsl2 的包
modprobe aufs

# AUFS 挂载
mount -t aufs -o dirs=./container-layer:./image-4:.image-3:./image-2:./image-1 none ./mnt

# 执行后文件被挂载到 ./mnt 里面, mount 的是目录，但是是目录里的文件被挂进去了。分层次的。
ls ./mnt -al

# 此时可以在文件系统 aufs 看到一条此挂载目录的信息 si_xxxx。
# 列举
ls /sys/fs/aufs

# 打印 aufs 配置
cat /sys/fs/aufs/config

# 打印信息
cat /sys/fs/aufs/<siid>/*

# 往 mnt image-1.txt 写入信息。
echo "append to image-1" >> ./mnt/image-1.txt

# 此时 原本 image-1/image-1.txt 不变
cat ./image-1/image-1.txt

# container-layer 目录多出修改文件 image-1.txt
ls ./container-layer
cat ./container-layer/image-1.txt
```


```bash
# vbox 虚拟机要使用共享需要安装插件
# host 是 Windows 下共享文件夹正常，粘贴板好像几个版本都是有问题的，不管是虚拟的是 windows 还是 ubuntu
apt install virtualbox-guest-utils

# ubuntu 安装 golang 也可以去官网下载 bin 版本解压配置路径，很麻烦，没下面简单。
apt  install golang-go
```