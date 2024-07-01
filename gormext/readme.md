# 转化 GORM model 到 TS 脚本

digdemo 里面有 gorm
直接搜索 gentool 查看其他。

```bash
gentool -db mysql -dsn "root:123456@tcp(localhost:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local" -tables "e_employee" -modelPkgName="models" -outPath="./entities" -fieldNullable -fieldWithIndexTag -fieldWithTypeTag  -fieldSignable 
```