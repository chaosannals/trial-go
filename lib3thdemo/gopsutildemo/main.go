package main

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	// 内存
	if v, err := mem.VirtualMemory(); err != nil {
		log.Fatalln(err)
	} else {
		// almost every return value is a struct
		fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

		// convert to JSON. String() is also implemented
		fmt.Println(v)
	}

	// CPU
	if v, err := cpu.Info(); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(v)
	}

	// CPU 占用，是否拆分
	if v, err := cpu.Times(false); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(v)
	}
}
