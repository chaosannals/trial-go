# dll demo

需要：

1. import "C"
2. //export 函数名
3. main 函数必须，但是可以放空。

```bat
@rem 执行编译，会生成 dlldemo.dll 和 dlldemo.h
go build -buildmode=c-shared -o dlldemo.dll main.go
```