package main

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
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

	fmt.Println("====================================")
	// 进程
	if ps, err := process.Processes(); err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(ps)
		for i, p := range ps {
			name, err := p.Name()
			if err != nil {
				fmt.Printf("error name: %d %v\n", p.Pid, err)
			}
			cpuPercent, err := p.CPUPercent()
			if err != nil && p.Pid > 4 {
				fmt.Printf("error cpu: %d %v\n", p.Pid, err)
			}
			fmt.Printf(
				"[%d] %s %f\n",
				i,
				name,
				cpuPercent,
			)
		}
	}
}
