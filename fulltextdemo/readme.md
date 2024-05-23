# 全文搜索 示例

## bleve demo

一个库。功能比较多。

```bash
# 命令行工具
go install github.com/blevesearch/bleve/v2/cmd/bleve@latest

# 查看
bleve --help

# 打印索引字段
bleve fields foo.bleve

# 查询
bleve query foo.bleve "中文" --highlight --fields -x
```

## gofound

这个是一个完整的程序，启动一个 HTTP 服务，向外提供接口。功能少，简单。

