package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	d := filepath.Dir(os.Args[0])
	if len(os.Args) > 1 {
		d = filepath.Join(d, os.Args[1])
	}
	r.StaticFS("/", http.Dir(d))
	r.Run(":50000")
}
