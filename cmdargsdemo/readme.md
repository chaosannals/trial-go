# 命令行库示例

## kingpin

面向过程风格。
可以单命令也可以子命令模式。

```bash
# -- 后是传递的默认参数
go run ./kingpindemo -- name

#  查看
go run ./kingpindemo --help
```

## go-flags

面向对象风格。
可以单命令也可以子命令模式。

这个库不同系统的命令会区分 Windows 是 /n, Linux 下是 -n.

```bat
@rem 多个参数需要重复 /ptrslice 类似这种，而不是紧随其后。
go run ./goflagsdemo /n name /ptrslice 123 /ptrslice 335 /ptrslice 1231
```

## kong

面向对象风格，多参数是紧随其后。
子命令模式。

```bash
#
go run ./kongdemo --help

# 多路径父级和当前
go run ./kongdemo ls .. .
```