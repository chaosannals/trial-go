package main

import (
	_ "image/gif"  // 允许加载的图片 GIF 格式库
	_ "image/jpeg" // 允许加载的图片 JPEG 格式库
	_ "image/png"  // 允许加载的图片 PNG 格式库

	"fmt"
	"time"
	"syscall/js"
	"github.com/xuri/excelize/v2"
	"github.com/chaosannals/trial-go-wasmxls/utils"
)

func init() {

}

func makeExcel(items js.Value) error {
	f := excelize.NewFile()
	si := f.NewSheet("Sheet")
	length := items.Length()

	for i := 0; i < length; i+=1 {
		fmt.Printf("写入项 %d", i)
		f.SetRowHeight("Sheet", i, 50)
		item := items.Index(i)
		j := i + 2

		itemNoCell := fmt.Sprintf("A%d", j)
		itemNo := item.Get("编号").String()
		f.SetCellValue("Sheet", itemNoCell, itemNo)

		itemPicCell := fmt.Sprintf("B%d", j)
		itemPic := item.Get("大图").String()
		b, err := utils.Download(itemPic)
		if err != nil {
			fmt.Printf("下载失败：%s\n", itemPic)
		} else {
			fmt.Printf("下载完成：%s\n", itemPic)
			if err := f.AddPictureFromBytes(
				"Sheet",
				itemPicCell,
				"{ \"autofit\": true, \"lock_aspect_ratio\": true }",
				"pic",
				".jpg",
				b); err != nil {
				fmt.Printf("图片添加失败 %s \n", itemPic)
				fmt.Println(err)
			}
		}
	}

	f.SetActiveSheet(si)
	b, err := f.WriteToBuffer()
	if err != nil {
		return err
	}
	ourl := utils.GoBufferToJsObjectURL(b)
	utils.MakeADownloadToBody(ourl, "result.xlsx")
	return nil
}

func main() {
	fmt.Println("开始生成")
	args := js.Global().Get("JS_TO_WASM_DATA")
	if !args.IsUndefined() {
		items := args.Get("清单数据")
		start := time.Now()
		makeExcel(items)
		end := time.Now()
		interval := end.Sub(start)
		fmt.Println(interval)
	}
}
