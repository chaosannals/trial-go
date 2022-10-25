# Welcome to Revel

A high-productivity web framework for the [Go language](http://www.golang.org/).

```bash
# 安装命令行工具
go install github.com/revel/cmd/revel@latest

# 创建项目
revel new myapp
```

### Start the web server:

```bash
# 在目录内
revel run

# 在目录外
revel run reveldemo
# 运行后在 网页的  /@tests 路由下有测试页面。
# 由配置的 module.testrunner 测试器执行， prod 模式 没有配置就没有该路由。

# 运行 命令行 测试  调整了几次参数还是会卡住。。。还是直接用 /@tests 测试就好。
revel test reveldemo [dev|prod] [AppTest|..]
```

### Go to http://localhost:9000/ and you'll see:

    "It works"

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


## Help

* The [Getting Started with Revel](http://revel.github.io/tutorial/gettingstarted.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/examples/index.html).
* The [API documentation](https://godoc.org/github.com/revel/revel).

