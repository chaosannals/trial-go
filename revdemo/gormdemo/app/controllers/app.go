package controllers

import (
	"fmt"
	"gormdemo/app"
	"gormdemo/util"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
	Count int
	IntPtr *int
}

func (c App) Index() revel.Result {
	*c.IntPtr += 10
	c.Count += 10
	c.ViewArgs["count"] = strconv.Itoa(c.Count)
	c.ViewArgs["pointer"] = unsafe.Pointer(&c)
	c.ViewArgs["intp"] = strconv.Itoa(*c.IntPtr)
	revel.AppLog.Info(fmt.Sprintf("intp: %d", *c.IntPtr))
	app.Log.Infof("---------------- %d", *c.IntPtr)
	return c.Render()
}

func (c App) ReturnMyResult() revel.Result {
	revel.AppLog.Warn("result call------")
	return &util.MyResult{
		Content: map[string]interface{}{
			"aaaaa": 123,
			"bbbb": 23423,
		},
	}
}

func (c App) GetType() reflect.Type {
	return reflect.TypeOf(c)
}

func (c App) GetValue() reflect.Value {
	return reflect.ValueOf(c)
}

func (c *App) GetPointerValue() reflect.Value {
	return reflect.ValueOf(c)
}

func (c *App) Inject(a int) {
	if c.IntPtr == nil {
		ip := 10
		c.IntPtr = &ip
	}
	c.Count += a
}