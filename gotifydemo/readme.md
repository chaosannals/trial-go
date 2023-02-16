# Gotify 消息推送服务

- gotify/server  服务器
- gotify/cli  命令行客户端

(gotify GitHub)[https://github.com/gotify] 各种语言 SDK

## 配置

- ./config.yml
- /etc/gotify/config.yml

```yml
server:
  keepaliveperiodseconds: 0 # keepalive 时间 0 = Go默认(15秒); -1 = 禁用;
  listenaddr: "" # 监听地址, 放空监听所有 0.0.0.0
  port: 80 # HTTP 端口
  ssl:
    enabled: false # true 启用 HTTPS
    redirecttohttps: true # 强制 HTTPS
    listenaddr: "" # 监听地址, 放空监听所有 0.0.0.0
    port: 443 # HTTPS 端口
    certfile: # 证书 (放空使用 letsencrypt)
    certkey: # 证书 key (放空使用 letsencrypt)
    letsencrypt:
      enabled: false # true 从 letsencrypt 申请
      accepttos: false # if you accept the tos from letsencrypt
      cache: data/certs # the directory of the cache from letsencrypt
      hosts: # 申请 letsencrypt 证书的域名列表
  #     - mydomain.tld
  #     - myotherdomain.tld
  responseheaders: # 响应头列表，默认空
  # X-Custom-Header: "custom value"

  cors: # 跨域设置
    alloworigins:
    # - ".+.example.com"
    # - "otherdomain.com"
    allowmethods:
    # - "GET"
    # - "POST"
    allowheaders:
  #   - "Authorization"
  #   - "content-type"

  stream:
    pingperiodseconds: 45 # the interval in which websocket pings will be sent. Only change this value if you know what you are doing.
    allowedorigins: # allowed origins for websocket connections (same origin is always allowed, default only same origin)
#     - ".+.example.com"
#     - "otherdomain.com"
database: # see below
  dialect: sqlite3
  connection: data/gotify.db

  # dialect: mysql
  # connection: root:password@/gotifydb?charset=utf8mb4&parseTime=True&loc=Local

  # dialect: postgres
  # connection: host=localhost port=3306 user=gotify dbname=gotify password=secret
defaultuser: # on database creation, gotify creates an admin user (these values will only be used for the first start, if you want to edit the user after the first start use the WebUI)
  name: admin # the username of the default user
  pass: admin # the password of the default user
passstrength: 10 # the bcrypt password strength (higher = better but also slower)
uploadedimagesdir: data/images # the directory for storing uploaded images
pluginsdir: data/plugins # the directory where plugin resides (leave empty to disable plugins)
registration: false # enable registrations
```

可以设置环境变量，但是不推荐，都是 yml 里面的。

```bash
GOTIFY_SERVER_PORT=80
GOTIFY_SERVER_KEEPALIVEPERIODSECONDS=0
GOTIFY_SERVER_LISTENADDR=
GOTIFY_SERVER_SSL_ENABLED=false
GOTIFY_SERVER_SSL_REDIRECTTOHTTPS=true
GOTIFY_SERVER_SSL_LISTENADDR=
GOTIFY_SERVER_SSL_PORT=443
GOTIFY_SERVER_SSL_CERTFILE=
GOTIFY_SERVER_SSL_CERTKEY=
GOTIFY_SERVER_SSL_LETSENCRYPT_ENABLED=false
GOTIFY_SERVER_SSL_LETSENCRYPT_ACCEPTTOS=false
GOTIFY_SERVER_SSL_LETSENCRYPT_CACHE=certs
# lists are a little weird but do-able (:
# GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS=- mydomain.tld\n- myotherdomain.tld
GOTIFY_SERVER_RESPONSEHEADERS="X-Custom-Header: \"custom value\""
# GOTIFY_SERVER_CORS_ALLOWORIGINS="- \".+.example.com\"\n- \"otherdomain.com\""
# GOTIFY_SERVER_CORS_ALLOWMETHODS="- \"GET\"\n- \"POST\""
# GOTIFY_SERVER_CORS_ALLOWHEADERS="- \"Authorization\"\n- \"content-type\""
# GOTIFY_SERVER_STREAM_ALLOWEDORIGINS="- \".+.example.com\"\n- \"otherdomain.com\""
GOTIFY_SERVER_STREAM_PINGPERIODSECONDS=45
GOTIFY_DATABASE_DIALECT=sqlite3
GOTIFY_DATABASE_CONNECTION=data/gotify.db
GOTIFY_DEFAULTUSER_NAME=admin
GOTIFY_DEFAULTUSER_PASS=admin
GOTIFY_PASSSTRENGTH=10
GOTIFY_UPLOADEDIMAGESDIR=data/images
GOTIFY_PLUGINSDIR=data/plugins
GOTIFY_REGISTRATION=false
```