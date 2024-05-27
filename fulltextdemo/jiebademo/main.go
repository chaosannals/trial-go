package main

import (
	"fmt"

	"github.com/wangbin/jiebago"
)

var seg jiebago.Segmenter

func init() {
	seg.LoadDictionary("dict.txt")
}

func print(ch <-chan string) {
	for word := range ch {
		fmt.Printf(" %s /", word)
	}
	fmt.Println()
}

func main() {
	fmt.Print("【全模式】：")
	print(seg.CutAll("我来到北京清华大学"))

	fmt.Print("【精确模式】：")
	print(seg.Cut("我来到北京清华大学", false))

	fmt.Print("【新词识别】：")
	print(seg.Cut("他来到了网易杭研大厦", true))

	fmt.Print("【搜索引擎模式】：")
	print(seg.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))

	fmt.Print("【搜索引擎模式】：")
	print(seg.CutForSearch("Hello world,小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))

	fmt.Print("【搜索引擎模式】：")
	print(seg.CutForSearch("You can you up, no can no BB.小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))
}
