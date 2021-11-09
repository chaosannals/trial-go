# [trial-go](https://github.com/chaosannals/trial-go)

## 命令

```bash
# 查看配置
go env
```

## Go 构建镜像

```sh
docker run -v /host/workspace:/workspace -e GOPROXY='https://mirrors.aliyun.com/goproxy/' -e GO111MODULE=on --name gomake gomake
```
