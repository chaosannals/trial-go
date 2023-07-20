# [go-zero](https://go-zero.dev/) Demo

```bash
# 安装脚手架工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 通过脚手架工具安装 protoc
goctl env check --install --verbose --force
```

VSCODE 安装插件，搜索 goctl 安装。支持期自定义的 *.api 文件语法高亮。

多个 api 文件会导致生成多个 main ，此框架应该只能单个 api 文件。

```bash
# 创建 API 项目
goctl api new apidemo

# 【项目目录执行】方式和一般 go 项目一致。
# 整理依赖文件
go mod tidy
# 启动 go 程序
go run apidemo.go

# 创建 api 文件
goctl api -o user.api

# api 文件生成 golang 代码，生成文件存在就不生成。
goctl api go --dir . --api apidemo.api
# 文件内容迁移到 apidemo.api 里
# goctl api go --dir . --api user.api

# api 生成文档
goctl api doc --dir . --o ./docs

# api 文件格式化当前目录所有 *.api 文件
goctl api format --dir .
```

```bash
# 创建 grpc 项目
goctl rpc new grpcdemo

# 创建 grpc 文件
goctl rpc --o greet.proto
```


```bash
# 生成 mysql 模型代码
# 这框架自己实现的 ORM 代码只是 mysql 预处理字符串模板，ORM 能力很弱。

# 通过 DDL 的 SQL 文件生成模型，算做了一半的 dbfirst 吧。
goctl model mysql ddl --src user.sql --dir .
```
