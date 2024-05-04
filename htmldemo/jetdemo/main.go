package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./jetviews"), // 相对 PWD 目录
	jet.InDevelopmentMode(),
)

func main() {
	views.AddGlobalFunc("base64", func(args jet.Arguments) reflect.Value {
		args.RequireNumOfArguments("base64", 1, 1)
		buffer := bytes.NewBuffer(nil)
		fmt.Fprint(buffer, args.Get(0))
		return reflect.ValueOf(base64.URLEncoding.EncodeToString(buffer.Bytes()))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view, err := views.GetTemplate("index.jet")
		if err != nil {
			log.Printf("模板错误 %v \n", err)
		}
		view.Execute(w, nil, map[string]any{
			"title": "标题",
		})
	})
	http.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		view, err := views.GetTemplate("page.jet")
		if err != nil {
			log.Printf("模板错误 %v \n", err)
		}
		view.Execute(w, nil, map[string]any{
			"title": "标题",
		})
	})

	http.ListenAndServe(":12345", nil)
}
