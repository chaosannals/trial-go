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
