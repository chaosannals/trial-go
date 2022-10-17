# google wire demo 依赖注入

```bash
# 全局安装命令行工具
go install github.com/google/wire/cmd/wire@latest

# 会找到当前目录的 wire.go 文件生成 wire_gen.go 文件
# 文件会替换 wire.Build 使用的函数，生成注入代码。
# 由于函数名重复，所以很不方便，这个还不如自己手写。
wire 

# 指定文件生成 wire_gen.go
wire my.go
```