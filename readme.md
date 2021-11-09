# [trial-go](https://github.com/chaosannals/trial-go)

## 命令

```bash
# 查看配置
go env
```

## 多版本

```bash
# 拉后会在 GOPATH 的 bin 目录下看到，有的会多层平台文件夹里面。
go get golang.org/dl/go1.17

# 然后通过该版本执行特定的命令前，必须先下载。
go1.17 download

# 之后就用版本号的版本执行命令。
# 带了平台的处理起来比较麻烦。
# 比如 win10 下 go1.10 就在 windows_386 目录下，是 32位的。
```

## Go 构建镜像

```sh
docker run -v /host/workspace:/workspace -e GOPROXY='https://mirrors.aliyun.com/goproxy/' -e GO111MODULE=on --name gomake gomake
```
