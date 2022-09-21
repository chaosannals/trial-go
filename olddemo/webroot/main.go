package main

import (
	"net/http"
	"os"
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	d := filepath.Dir(os.Args[0])
	if len(os.Args) > 1 {
		d = filepath.Join(d, os.Args[1])
	}
	fmt.Println("http://localhost:50000")
	r.StaticFS("/", http.Dir(d))
	r.Run(":50000")
}
