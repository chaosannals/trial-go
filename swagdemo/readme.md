# swag demo

```bash
# 安装工具
go install github.com/swaggo/swag/cmd/swag@latest

# 根据注解 生成 docs 目录
# 有些版本生成的代码有点问题，会多出字段，删了就 OK 。
swag init

# 不同框架使用不同设配库。
# gin 的设配
go get -u github.com/swaggo/gin-swagger

# echo 的设配
go get -u github.com/swaggo/echo-swagger
```
