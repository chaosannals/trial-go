package main

import (
	_ "gfdemo/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"gfdemo/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
