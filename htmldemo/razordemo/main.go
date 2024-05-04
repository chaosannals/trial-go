package main

import (
	"fmt"
	"net/http"

	"github.com/chaosannals/htmldemo/razordemo/models"
	"github.com/chaosannals/htmldemo/razordemo/tpl"
)

func main() {
	user := &models.User{}
	user.Name = "go"
	user.Email = "hello@world.com"
	user.Intro = "I love gorazor!"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, tpl.Home(1, user))
	})

	http.ListenAndServe(":12345", nil)
}
