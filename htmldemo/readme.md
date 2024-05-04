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