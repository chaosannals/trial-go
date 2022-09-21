package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"github.com/gonutz/w32"
)

func getVersion(path string) string {
	size := w32.GetFileVersionInfoSize(path)
	if size > 0 {
		info := make([]byte, size)
		ok := w32.GetFileVersionInfo(path, info)
		if !ok {
			fmt.Println("GetFileVersionInfo faield")
			os.Exit(-5)
		}
		fixed, ok := w32.VerQueryValueRoot(info)
		if !ok {
			fmt.Println("VerQueryValueRoot faield")
			os.Exit(-5)
		}
		version := fixed.FileVersion()
		return fmt.Sprintf(
			".%d.%d.%d.%d",
			version&0xFFFF000000000000>>48,
			version&0x0000FFFF00000000>>32,
			version&0x00000000FFFF0000>>16,
			version&0x000000000000FFFF>>0,
		)
	}
	return ""
}

func main() {
	fmt.Println(filepath.Base(os.Args[0]))
	name := path.Base(filepath.ToSlash(os.Args[0]))
	fmt.Println(name)
	version := getVersion(os.Args[0])
	fmt.Println(version)
}