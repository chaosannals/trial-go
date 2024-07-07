package main

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
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

	fmt.Println("====================================")
	// CPU
	if v, err := cpu.Info(); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(v)
	}

	fmt.Println("====================================")
	// CPU 占用，是否拆分
	if v, err := cpu.Times(false); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(v)
	}

	fmt.Println("====================================")
	// 硬盘分区
	if dps, err := disk.Partitions(true); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(dps)
		for i, dp := range dps {
			// Windows 下调用提示没有实现。not implemented yet
			// if sn, err := disk.SerialNumber(dp.Mountpoint); err != nil {
			// 	log.Fatalln(err)
			// } else {
			// 	fmt.Printf("[%d] %v\n", i, sn)
			// }
			// Windows 下调用提示没有实现。not implemented yet
			// if dl, err := disk.Label(dp.Mountpoint); err != nil {
			// 	log.Fatalln(err)
			// } else {
			// 	fmt.Printf("[%d] %v\n", i, dl)
			// }
			if du, err := disk.Usage(dp.Mountpoint); err != nil {
				log.Fatalln(err)
			} else {
				fmt.Printf("[%d] %v\n", i, du)
			}
		}
	}
}
