# ent demo

好像所有工具都加了 -mod 导致不能再 work 模式下使用。
而且还生成巨多文件。。

```bash
# 安装 cmd 工具

# readme 全局安装了。按理说与下面命令等价。
go install entgo.io/ent/cmd/ent@latest
# 但是还是需要执行这个命令才能使用。
go get -d entgo.io/ent/cmd/ent
```

```bash
# 生成命令，生成模板 codeFirst

# mod 下执行
# 执行这个命令后要自己去修改文件添加字段，不然只有一个 id 字段。
# 不是 dbFirst ，没有根据数据库生成。
# 只是 codeFirst 的生成模板命令
# ./ent/schema 里面会多出 user  pet 的定义
go run -mod=mod entgo.io/ent/cmd/ent new User Pet

# work 下去掉 mod (此命令可以执行，但是 ent 目前看并不支持 work)
go run entgo.io/ent/cmd/ent new User Pet

# 生成代码文件 codeFirst 只是根据 User Pet 模块生成 DAO 函数
# 修改定义后需要重新执行，生成新代码。
go generate ./ent
```


```bash
# 查看表结构
go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema
```