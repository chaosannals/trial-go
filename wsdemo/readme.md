# Linux Namespace Demo

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