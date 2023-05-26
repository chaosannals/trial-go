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
# 启动一个 stress 进程 200M 内存
stress --vm-bytes 200m --vm-keep -m 1

# 在 /sys/fs/cgroup/memory 子系统里面建子 cgroup 并进入
mkdir test-limit-memory
cd test-limit-memory

# 写入内存限制
sudo sh -c "echo '100m' > memory.limit_in_bytes"

# 把当前进程挂到 cgroup 上去，此时当前的 sh 被新的 test-limit-memory 限制内存。
sudo sh -c "echo $$ > tasks" 
sudo sh -c "echo $$ > cgroup.procs" 
```
