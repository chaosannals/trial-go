# wstpd

基于 websocket 的 软件中间传输服务（前端 到 软件端）。

```bat
@rem 设置阿里云代理
set GOPROXY=https://mirrors.aliyun.com/goproxy/
go mod tidy
```

```bash
# 设置阿里云代理
export GOPROXY=https://mirrors.aliyun.com/goproxy/
go mod tidy
```

## 架构说明

前端通过 websocket 连接 wstpd 服务器，每次调用都需要传递 appkey 和 token （由后端生成给到前端用来做认证）。传递的任务被转发到软件，单次请求可能获得多次响应，因为过程很漫长，响应会返回进度。

软件通过 websocket 连接 wstpd 服务器，每分钟执行一次心跳连接保证服务的连接，如果连接断开重试连接。等待 wstpd 转发过来的前端任务，完成部分任务便返回进度，最终完成任务要回馈完成。

注：由于跨域且高版本浏览器必须是 wss 协议，所以数据的传输过程本身就是加密的。

## 

## 开发与部署

### 构建

```bash
# 构建并上传阿里云镜像仓库
./build.ps1
```

### Nginx 配置

```nginx
upstream wstpd {
    server 127.0.0.1:44444;
}

server {
    # .......

    location /svc {
        proxy_http_version 1.1;
        proxy_set_header Upgrade websocket;
        proxy_set_header Connection "Upgrade";
        proxy_read_timeout 600s;
            
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $http_host;
        proxy_pass http://wstpd;
    }
    
    location /biz {
        proxy_http_version 1.1;
        proxy_set_header Upgrade websocket;
        proxy_set_header Connection "Upgrade";
        proxy_read_timeout 600s;

        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $http_host;
        proxy_pass http://wstpd;
    }
}
```
