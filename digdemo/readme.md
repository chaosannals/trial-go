# uber dig demo 依赖注入

## gorm CodeFist



## gorm DbFirst

```bash
# 安装工具
go install gorm.io/gen/tools/gentool@latest

# 工具说明
gentool -h

# 生成
# -outPath="./migrations" query 和 codefirst 文件路径
# -modelPkgName="models" 模型的包名，一般和 query 等是分开的、
# -withUnitTest 输出单元测试
# -onlyModel 只输出模型
gentool -db mysql -dsn "root:password@tcp(localhost:3306)/exert?charset=utf8mb4&parseTime=True&loc=Local" -tables "e_employee,e_employee_mobilephone" -modelPkgName="models" -outPath="./entities" -fieldNullable -fieldWithIndexTag -fieldWithTypeTag  -fieldSignable 

gentool -dsn "root:password@tcp(localhost:3306)/exert?charset=utf8mb4&parseTime=True&loc=Local" -tables "e_employee,e_employee_mobilephone" -modelPkgName="models" -outPath="./entities" -withUnitTest -fieldWithIndexTag
```

```bash
# 生成 query 和 codefirst 依赖
go get gorm.io/gen


# 生成的单元测试依赖这个驱动
go get gorm.io/driver/sqlite

```