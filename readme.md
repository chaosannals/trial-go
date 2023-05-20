# [trial-go](https://github.com/chaosannals/trial-go)

## 命令

```bash
# 查看配置
go env

# go1.13 后 可以通过命令设置
go env -w GOPROXY=https://goproxy.io,direct
go env -w GO111MODULE=on
```

### mod 模块

```bash
# 开启 mod
go env -w GO111MODULE=on

# 初始化一个模块
go mod init github.com/chaosannals/project

# 清理 go.mod 依赖
go mod tidy
```

### work 工作区

```bash
# 初始化一个工作区在当前（空）
go work init

# 指定第一个 mod 目录 需要先用 mod 命令初始化
go work init ./moddir

# 添加 mod 进入 work
go work use ./moddir
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

```bash
# 使用 go env -w 设置
go env -w GOPROXY=https://proxy.golang.com.cn,direct
# 不用像下面那样改环境变量。

# 另一个镜像
go env -w GOPROXY=https://goproxy.cn,direct
```

```bat
@rem 设置阿里镜像
set GOPROXY=https://mirrors.aliyun.com/goproxy/
```

```bash
# 设置阿里镜像
export GOPROXY=https://mirrors.aliyun.com/goproxy/
```

```sh
docker run -v /host/workspace:/workspace -e GOPROXY='https://mirrors.aliyun.com/goproxy/' -e GO111MODULE=on --name gomake gomake
```


## build

### 交叉编译

```bash
# 设置环境变量
env GOOS=linux GOARCH=amd64

# 构建
go build -o target
```

注：windows 下要使用 cmd 而不是 pwsh 执行，且set 变量时 && 前面的空格会被弄到变量值里面，所以不能要空格，要紧挨着。
```cmd
set CGO_ENABLED=0&&set GOOS=js&&set GOARCH=wasm&&go build -o game.wasm main.go
```
指定生成 wasm 文件。


```bash
# 查看支持构建的种类
go tool dist list
```
