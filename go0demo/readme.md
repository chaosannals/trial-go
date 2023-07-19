# [go-zero](https://go-zero.dev/) Demo

```bash
# 安装脚手架工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 通过脚手架工具安装 protoc
goctl env check --install --verbose --force
```

VSCODE 安装插件，搜索 goctl 安装。支持期自定义的 *.api 文件语法高亮。

```bash
# 创建 API 项目
goctl api new apidemo

# 【项目目录执行】方式和一般 go 项目一致。
# 整理依赖文件
go mod tidy
# 启动 go 程序
go run apidemo.go

# api 文件生成 golang 代码，生成文件存在就不生成。
goctl api go --dir . --api apidemo.api

# api 生成文档
goctl api doc --dir . --o ./docs

# api 文件格式化当前目录所有 *.api 文件
goctl api format --dir .
```
