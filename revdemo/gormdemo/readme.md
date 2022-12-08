# REVEL GORM DEMO

## Revel

```bash
revel run
# /@tests 下是接口测试
```

## 两种数据库管理的方式 GEN 和 GEN TOOL

### [GEN](https://gorm.io/gen/index.html)

```bash
# 通过定制化程序生成
go run tool/main.go
```

### [GEN Tool](https://gorm.io/gen/gen_tool.html)

这个是官方封装的一个命令行程序，通过配置参数，不用写 go 文件。

注：没有开放 Mode 参数，导致 Mode 是空配置。

```bash
# 安装 GORM 工具
go install gorm.io/gen/tools/gentool@latest

# 工具说明
gentool -h

# 生成
# -outPath="./migrations" query 和 codefirst 文件路径
# -modelPkgName="models" 模型的包名，一般和 query 等是分开的、
# -withUnitTest 输出单元测试
# -onlyModel 只输出模型
gentool -db mysql -dsn "root:password@tcp(localhost:3306)/exert?charset=utf8mb4&parseTime=True&loc=Local" -tables "e_employee,e_employee_mobilephone" -modelPkgName="models" -outPath="./entities" -fieldNullable -fieldWithIndexTag -fieldWithTypeTag  -fieldSignable 
```

注：digdemo 也有 gorm 的示例，不过是基于 echo + dig
