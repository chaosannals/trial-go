package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/chaosannals/swagdemo/echodemo/docs"
)

func main() {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := 1323
	fmt.Printf("http://127.0.0.1:%d/swagger/\n", port)

	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
