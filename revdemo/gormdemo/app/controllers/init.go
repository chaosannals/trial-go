package controllers

import (
	"gormdemo/util"

	"github.com/revel/revel"
)

func doThing(c *revel.Controller) revel.Result {
	revel.AppLog.Warnf("result: %v", c.Result)
	if r, ok := c.Result.(*util.MyResult); ok {
		return c.RenderJSON(map[string]interface{}{
			"content": r.Content,
			"tip":     "doNothing",
		})
	}
	return nil
}

func init() {
	revel.InterceptFunc(doThing, revel.AFTER, &App{})
}
