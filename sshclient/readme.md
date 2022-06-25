# ssh client

```ini
ssh_host=127.0.0.1
ssh_port=22
ssh_user=root
ssh_password=123456789
ssh_key_path=
```

```bash
# 打包，v2 版本使用下面安装 cmd 工具
go install fyne.io/fyne/v2/cmd/fyne@latest

# 打包指定平台
fyne package -os darwin -icon icon.png
fyne package -os linux -icon icon.png
fyne package -os windows -icon icon.png
fyne package -os android -appID com.example.demo -icon icon.png
fyne package -os ios - appID com.example.demo -icon icon.png
```

注：windows 下 Msys2 TDM-gcc cygwin 任选一，TDM-gcc 比较简单，装上即可。
有些软件可能会不知道什么时候给你装上，所以如果没有缺少gcc的编译提示，就不用了。