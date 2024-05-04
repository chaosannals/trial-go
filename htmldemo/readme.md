# html 模板引擎示例

## goview

这个输出是指定死了 httpResponseWriter 的。
语法上基本和其他模板引擎一致。

没有看到 循环语句等结构化语句。

TODO 没有看到 字典 多级变量的模板语法。

## jet

语法比较多。
定义函数比较费劲，没有 goview 直观。除了这点应该全面优于 goview
输出接口是 io.Writer 比较好。goview 这种 http ResponseWriter 就很坑。
耦合比较少，没有约定目录结构，只指定 Root ，就比较自由。
结构化语法齐全 range

import 使用 block 管理，比其他引擎优秀。

### 模板语法

- {* 注释 *}
- {{ v := "text" }}  定义变量
- {{ s := slice("aaa", "bbb", "ccc")}} 数组或切片
- {{ s[1] }}
- {{ m := map("aaa", 1, "bbb", 2)}} 字典
- {{ m["aaa"] }}
- {{ .HasTitle ? .Title : "default title" }} 三元运算符

```html
{{ yield body() }}
<!-- 作为布局而被扩展的模板，有个固定名为 "body" 的槽位指定 继承模板 的内容块 -->
```

```html
{{ extends "layouts/simple.jet" }}
```

```html
{{ import "partials/ad.jet"}} <!-- 必须放在文件头部 -->

{{range users}}
    <div>{{.Name}}</div>
{{end}}

{{ yield body() }}
```

## gorazor

这个是把 C# 的 Razor 模板语言移植到 go 了。

搞得和 C# partial 类一个搞法，限制太死。每个模板文件得配一个 go 文件。

采用生成的方式使用，通过 gohtml 文件生成一个 go 文件。

有个非常不好的点，layout 必须包名叫 layout ，约定了这个固定名称。

这种会导致把逻辑写在模板文件里面。和 razor 一个德性。而且golang 没有 partial class ，所以代码无法分离。


```bash
# 这个生成器有 bug ，官方示例生成模板有问题。最后直接复制的官方生成物。
go install github.com/sipin/gorazor@latest
go install github.com/sipin/gorazor@1.2.2
go install github.com/sipin/gorazor@1.2.1
go install github.com/sipin/gorazor@1.0
```

```bash
#Process folder: 
gorazor template_folder output_folder

#Process file: 
gorazor template_file output_file
```