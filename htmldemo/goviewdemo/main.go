package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/foolin/goview"
)

// 模板用法：
// 变量  {{.name}}
// 定义槽位  {{template "name" .}}
// 填充槽位 {{define "name"}} 内容 {{end}}
// 引入模板 {{include "name"}}
func main() {
	gv := goview.New(goview.Config{
		Root:      "goviews", // 相对 PWD 路径
		Extension: ".tpl",
		Master:    "layouts/master", // master.tpl 这里省略了后缀
		Partials: []string{ // 全局的部件模板 ad.tpl 省略了后缀
			"partials/ad",
		},
		Funcs: template.FuncMap{ // 自定义全局模板函数
			"sub": func(a, b int) int {
				return a - b
			},
			"year": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
		Delims:       goview.Delims{Left: "{{", Right: "}}"}, // 模板的标识符定义
	})
	goview.Use(gv)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// index 是模板名 index.tpl 省略了后缀，此路径相对 Root
		err := goview.Render(w, http.StatusOK, "index", goview.M{
			"title": "Index title",
			"add": func(a, b int) int {
				return a + b
			},
		})
		if err != nil {
			fmt.Fprintf(w, "模板引擎错误 %v", err)
		}
	})

	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		// 不应该在模板名上加后缀。
		err := goview.Render(w, http.StatusOK, "page", goview.M{
			"title": "Page",
		})
		if err != nil {
			fmt.Fprintf(w, "page 模板引擎错误 %v", err)
		}
	})

	http.ListenAndServe(":12345", nil)
}
