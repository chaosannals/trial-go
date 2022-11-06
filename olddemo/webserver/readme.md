# webserver

Windows 下更新了 SQLITE(github.com/mattn/go-sqlite3) 库后就报警告：
```bash
# issue 提议是配置这个
go env -w CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
# 在编译时自己加临时后报错，感觉不可用。
set CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
```